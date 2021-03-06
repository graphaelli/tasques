package recurring

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/rs/zerolog/log"

	"github.com/lloydmeta/tasques/internal/domain/metadata"
	"github.com/lloydmeta/tasques/internal/domain/queue"
	"github.com/lloydmeta/tasques/internal/domain/task"
	"github.com/lloydmeta/tasques/internal/domain/task/recurring"
	"github.com/lloydmeta/tasques/internal/infra/elasticsearch/common"
)

var TasquesRecurringTasksIndex = ".tasques_recurring_tasks"

type EsService struct {
	client         *elasticsearch.Client
	scrollPageSize uint
	scrollTtl      time.Duration
	getUTC         func() time.Time // for mocking
}

func (e *EsService) SetUTCGetter(getter func() time.Time) {
	e.getUTC = getter
}

func NewService(client *elasticsearch.Client, scrollPageSize uint, scrollTtl time.Duration) recurring.Service {
	return &EsService{
		client:         client,
		scrollPageSize: scrollPageSize,
		scrollTtl:      scrollTtl,
		getUTC: func() time.Time {
			return time.Now().UTC()
		},
	}
}

func (e *EsService) Create(ctx context.Context, task *recurring.NewTask) (*recurring.Task, error) {
	toPersist := e.newToPersistable(task)
	toPersistBytes, err := json.Marshal(toPersist)
	if err != nil {
		return nil, common.JsonSerdesErr{Underlying: []error{err}}
	}
	req := esapi.CreateRequest{
		Index:      TasquesRecurringTasksIndex,
		DocumentID: string(task.ID),
		Body:       bytes.NewReader(toPersistBytes),
	}

	rawResp, err := req.Do(ctx, e.client)
	if err != nil {
		return nil, common.ElasticsearchErr{Underlying: err}
	}
	defer rawResp.Body.Close()
	statusCode := rawResp.StatusCode
	switch {
	case 200 <= statusCode && statusCode <= 299:
		var response common.EsCreateResponse
		if err := json.NewDecoder(rawResp.Body).Decode(&response); err != nil {
			return nil, common.JsonSerdesErr{Underlying: []error{err}}
		}
		domainTask := persistedToDomain(task.ID, &toPersist, response.Version())
		return &domainTask, nil
	case statusCode == 409:
		// do a get and see if it exists as a soft-deleted record
		existing, err := e.Get(ctx, task.ID, true)
		if err == nil && bool(existing.IsDeleted) {
			// Update (reindex) if the doc was *soft* deleted
			toPersistBytes, err := json.Marshal(toPersist)
			if err != nil {
				return nil, common.JsonSerdesErr{Underlying: []error{err}}
			}
			req := esapi.IndexRequest{
				Index:      TasquesRecurringTasksIndex,
				DocumentID: string(task.ID),
				Body:       bytes.NewReader(toPersistBytes),
			}
			rawResp, err := req.Do(ctx, e.client)
			if err != nil {
				return nil, common.ElasticsearchErr{Underlying: err}
			}
			defer rawResp.Body.Close()
			statusCode := rawResp.StatusCode
			switch {
			case 200 <= statusCode && statusCode <= 299:
				var response common.EsCreateResponse
				if err := json.NewDecoder(rawResp.Body).Decode(&response); err != nil {
					return nil, common.JsonSerdesErr{Underlying: []error{err}}
				}
				domainTask := persistedToDomain(task.ID, &toPersist, response.Version())
				return &domainTask, nil
			default:
				return nil, common.UnexpectedEsStatusError(rawResp)
			}
		} else {
			if err != nil {
				return nil, err
			} else {
				return nil, recurring.AlreadyExists{ID: task.ID}
			}
		}
	default:
		return nil, common.UnexpectedEsStatusError(rawResp)
	}
}

func (e *EsService) Get(ctx context.Context, id task.RecurringTaskId, includeSoftDeleted bool) (*recurring.Task, error) {
	req := esapi.GetRequest{
		Index:      TasquesRecurringTasksIndex,
		DocumentID: string(id),
	}
	rawResp, err := req.Do(ctx, e.client)
	if err != nil {
		return nil, common.ElasticsearchErr{Underlying: err}
	}
	defer rawResp.Body.Close()

	switch rawResp.StatusCode {
	case 200:
		var resp esHitPersistedRecurringTask
		if err := json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
			return nil, common.JsonSerdesErr{Underlying: []error{err}}
		}
		domainModel := resp.toDomainTask()
		if bool(domainModel.IsDeleted) && !includeSoftDeleted {
			return nil, recurring.NotFound{
				ID: id,
			}
		}
		return &domainModel, nil
	case 404:
		return nil, recurring.NotFound{
			ID: id,
		}
	default:
		return nil, common.UnexpectedEsStatusError(rawResp)
	}
}

func (e *EsService) Delete(ctx context.Context, id task.RecurringTaskId) (*recurring.Task, error) {
	existing, err := e.Get(ctx, id, false)
	if err != nil {
		return nil, err
	} else {
		existing.IntoDeleted()
		_, updateErr := e.Update(ctx, existing)
		if updateErr != nil {
			return nil, updateErr
		} else {
			return existing, nil
		}
	}
}

func (e *EsService) All(ctx context.Context) ([]recurring.Task, error) {
	searchBody := buildUndeletedListSearchBody(e.scrollPageSize)
	var found []recurring.Task
	err := e.refreshIndex(ctx)
	if err != nil {
		return nil, err
	}
	err = e.scanRecurringTasks(ctx, searchBody, e.scrollTtl, func(recurringTasks []recurring.Task) error {
		found = append(found, recurringTasks...)
		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return found, nil
	}
}

func (e *EsService) NotLoaded(ctx context.Context) ([]recurring.Task, error) {
	searchBody := buildNotLoadedSinceSearchBody(e.scrollPageSize)
	var found []recurring.Task
	err := e.scanRecurringTasks(ctx, searchBody, e.scrollTtl, func(recurringTasks []recurring.Task) error {
		found = append(found, recurringTasks...)
		return nil
	})
	if err != nil {
		return nil, err
	} else {
		return found, nil
	}
}

func (e *EsService) Update(ctx context.Context, update *recurring.Task) (*recurring.Task, error) {
	now := e.getUTC()
	update.LoadedAt = nil
	toPersist := domainToPersistable(update, metadata.ModifiedAt(now))
	toPersistBytes, err := json.Marshal(toPersist)
	if err != nil {
		return nil, common.JsonSerdesErr{Underlying: []error{err}}
	}
	req := esapi.IndexRequest{
		Index:         TasquesRecurringTasksIndex,
		DocumentID:    string(update.ID),
		Body:          bytes.NewReader(toPersistBytes),
		IfPrimaryTerm: esapi.IntPtr(int(update.Metadata.Version.PrimaryTerm)),
		IfSeqNo:       esapi.IntPtr(int(update.Metadata.Version.SeqNum)),
	}
	rawResp, err := req.Do(ctx, e.client)
	if err != nil {
		return nil, common.ElasticsearchErr{Underlying: err}
	}
	defer rawResp.Body.Close()
	respStatus := rawResp.StatusCode
	switch {
	case 200 <= respStatus && respStatus <= 299:
		// Updated, grab new metadata
		var resp common.EsUpdateResponse
		if err := json.NewDecoder(rawResp.Body).Decode(&resp); err != nil {
			return nil, common.JsonSerdesErr{Underlying: []error{err}}
		}
		update.Metadata.Version = resp.Version()
		return update, nil
	case respStatus == 409:
		return nil, recurring.InvalidVersion{ID: update.ID}
	case respStatus == 404:
		return nil, recurring.NotFound{ID: update.ID}
	default:
		return nil, common.UnexpectedEsStatusError(rawResp)
	}
}

func (e *EsService) MarkLoaded(ctx context.Context, toMarks []recurring.Task) (*recurring.MultiUpdateResult, error) {
	now := e.getUTC()
	loadedAt := recurring.LoadedAt(now)
	modifiedAt := metadata.ModifiedAt(now)
	for i := 0; i < len(toMarks); i++ {
		toMark := &toMarks[i]
		toMark.LoadedAt = &loadedAt
	}
	bulkReqBody, err := buildTasksBulkUpdateNdJsonBytes(toMarks, modifiedAt)
	if err != nil {
		return nil, err
	}
	claimBulkReq := esapi.BulkRequest{
		Body: bytes.NewReader(bulkReqBody),
	}
	rawResp, err := claimBulkReq.Do(ctx, e.client)
	if err != nil {
		return nil, common.ElasticsearchErr{Underlying: err}
	}
	defer rawResp.Body.Close()
	if rawResp.IsError() {
		return nil, common.UnexpectedEsStatusError(rawResp)
	}
	var response common.EsBulkResponse
	if err := json.NewDecoder(rawResp.Body).Decode(&response); err != nil {
		return nil, common.JsonSerdesErr{Underlying: []error{err}}
	}
	var multiResult recurring.MultiUpdateResult
	for updateTargetIdx, updateTarget := range toMarks {
		// This is guaranteed by ES
		result := response.Items[updateTargetIdx]
		resultInfo := result.Info()
		if resultInfo.IsOk() {
			updateTarget.Metadata.Version = resultInfo.Version()
			multiResult.Successes = append(multiResult.Successes, updateTarget)
		} else if resultInfo.Status == 404 {
			multiResult.NotFounds = append(multiResult.NotFounds, updateTarget)
		} else if resultInfo.Status == 409 {
			multiResult.VersionConflicts = append(multiResult.VersionConflicts, updateTarget)
		} else {
			otherError := recurring.BulkUpdateOtherError{
				RecurringTask: updateTarget,
				Result:        resultInfo.Result,
			}
			multiResult.Others = append(multiResult.Others, otherError)
		}
	}
	return &multiResult, nil
}

func (e *EsService) refreshIndex(ctx context.Context) error {
	req := esapi.IndicesRefreshRequest{
		Index:             []string{TasquesRecurringTasksIndex},
		AllowNoIndices:    esapi.BoolPtr(true),
		IgnoreUnavailable: esapi.BoolPtr(true),
	}
	rawResp, err := req.Do(ctx, e.client)
	if err != nil {
		return common.ElasticsearchErr{Underlying: err}
	}
	defer rawResp.Body.Close()
	respStatus := rawResp.StatusCode
	switch {
	case 200 <= respStatus && respStatus <= 299:
		return nil
	default:
		return common.UnexpectedEsStatusError(rawResp)
	}
}

// This is the main method that should be used for listing and scrolling through a potentially large collection of
// RecurringTasks
func (e *EsService) scanRecurringTasks(ctx context.Context, searchBody jsonObjMap, scrollTtl time.Duration, doWithBatch func(recurrings []recurring.Task) error) error {
	log.Debug().Msg("Beginning to scan recurring tasks")
	log.Debug().Interface("searchBody", searchBody).Msg("Scanning recurring tasks")
	recurringTasksWithScrollId, err := e.initSearch(ctx, searchBody, scrollTtl)
	if err != nil {
		return err
	}
	scannedRecurringTasks := recurringTasksWithScrollId.RecurringTasks
	var scrollIds []string
	scrollId := recurringTasksWithScrollId.ScrollId
	scrollIds = append(scrollIds, scrollId)
	defer func() {
		if scrollErr := e.clearScroll(ctx, scrollIds); scrollErr != nil && err == nil {
			err = scrollErr
		}
	}()

	for len(scannedRecurringTasks) > 0 {
		if err := doWithBatch(scannedRecurringTasks); err != nil {
			return err
		}
		nextTasksWithScrollId, err := e.scroll(ctx, scrollId, scrollTtl)
		if err != nil {
			return err
		}
		scannedRecurringTasks = nextTasksWithScrollId.RecurringTasks
		scrollId = nextTasksWithScrollId.ScrollId
		scrollIds = append(scrollIds, nextTasksWithScrollId.ScrollId)
	}
	log.Debug().Msg("Scanning recurring tasks end ")
	return nil
}

func (e *EsService) initSearch(ctx context.Context, searchBody jsonObjMap, scrollTtl time.Duration) (*recurringTasksWithScrollId, error) {
	searchBodyBytes, err := json.Marshal(searchBody)
	if err != nil {
		return nil, common.JsonSerdesErr{Underlying: []error{err}}
	}
	searchReq := esapi.SearchRequest{
		Scroll:            scrollTtl, // make this configurable
		Index:             []string{TasquesRecurringTasksIndex},
		AllowNoIndices:    esapi.BoolPtr(true),
		IgnoreUnavailable: esapi.BoolPtr(true), // Recurring Tasks Index might not exist yet
		Body:              bytes.NewReader(searchBodyBytes),
	}

	rawResp, err := searchReq.Do(ctx, e.client)
	if err != nil {
		return nil, common.ElasticsearchErr{Underlying: err}
	}
	defer rawResp.Body.Close()
	return processScrollResp(rawResp)
}

func (e *EsService) scroll(ctx context.Context, scrollId string, scrollTtl time.Duration) (*recurringTasksWithScrollId, error) {
	scrollReq := esapi.ScrollRequest{
		Scroll:   scrollTtl,
		ScrollID: scrollId,
	}

	rawResp, err := scrollReq.Do(ctx, e.client)
	if err != nil {
		return nil, common.ElasticsearchErr{Underlying: err}
	}
	defer rawResp.Body.Close()
	return processScrollResp(rawResp)
}

func processScrollResp(rawResp *esapi.Response) (*recurringTasksWithScrollId, error) {
	switch rawResp.StatusCode {
	case 200:
		var scrollResp esSearchScrollingResponse
		if err := json.NewDecoder(rawResp.Body).Decode(&scrollResp); err != nil {
			return nil, common.JsonSerdesErr{Underlying: []error{err}}
		}
		tasks := make([]recurring.Task, 0, len(scrollResp.Hits.Hits))
		for _, pTask := range scrollResp.Hits.Hits {
			tasks = append(tasks, pTask.toDomainTask())
		}
		return &recurringTasksWithScrollId{
			ScrollId:       scrollResp.ScrollId,
			RecurringTasks: tasks,
		}, nil
	default:
		return nil, common.UnexpectedEsStatusError(rawResp)
	}
}

func (e *EsService) clearScroll(ctx context.Context, scrollIds []string) error {
	if len(scrollIds) > 0 {
		clearScrollReq := esapi.ClearScrollRequest{ScrollID: scrollIds}
		rawResp, err := clearScrollReq.Do(ctx, e.client)
		if err != nil {
			return err
		} else {
			defer rawResp.Body.Close()
			switch rawResp.StatusCode {
			case 200:
				return nil
			default:
				return common.UnexpectedEsStatusError(rawResp)
			}
		}
	} else {
		return nil
	}
}

func buildUndeletedListSearchBody(pageSize uint) jsonObjMap {
	return jsonObjMap{
		"size":                pageSize,
		"seq_no_primary_term": true,
		"sort": []jsonObjMap{
			{
				"_id": jsonObjMap{
					"order": "asc",
				},
			},
		},
		"query": jsonObjMap{
			"bool": jsonObjMap{
				"filter": jsonObjMap{
					"bool": jsonObjMap{
						"must": []jsonObjMap{{
							"term": jsonObjMap{
								"is_deleted": false,
							},
						}},
					},
				},
			},
		},
	}
}

func buildNotLoadedSinceSearchBody(pageSize uint) jsonObjMap {
	return jsonObjMap{
		"size":                pageSize,
		"seq_no_primary_term": true,
		"sort": []jsonObjMap{
			{
				"_id": jsonObjMap{
					"order": "asc",
				},
			},
			{
				"metadata.modified_at": jsonObjMap{
					"order": "asc",
				},
			},
		},
		"query": jsonObjMap{
			"bool": jsonObjMap{
				"filter": jsonObjMap{
					"bool": jsonObjMap{
						"must": []jsonObjMap{
							{
								"bool": jsonObjMap{
									"must_not": []jsonObjMap{
										{
											"exists": jsonObjMap{
												"field": "loaded_at",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildTasksBulkUpdateNdJsonBytes(recurringTasks []recurring.Task, at metadata.ModifiedAt) ([]byte, error) {
	var errAcc []error
	var bytesAcc []byte
	for _, t := range recurringTasks {
		pair := buildUpdateBulkOp(&t, at)
		opBytes, err := json.Marshal(pair.op)
		if err != nil {
			errAcc = append(errAcc, err)
		}
		if len(errAcc) == 0 {
			bytesAcc = append(bytesAcc, opBytes...)
			bytesAcc = append(bytesAcc, "\n"...)
		}

		dataBytes, err := json.Marshal(pair.doc)
		if err != nil {
			errAcc = append(errAcc, err)
		}
		if len(errAcc) == 0 {
			bytesAcc = append(bytesAcc, dataBytes...)
			bytesAcc = append(bytesAcc, "\n"...)
		}
	}
	if len(errAcc) != 0 {
		return nil, common.JsonSerdesErr{Underlying: errAcc}
	} else {
		return bytesAcc, nil
	}
}

func buildUpdateBulkOp(task *recurring.Task, at metadata.ModifiedAt) updateRecurringTaskBulkOpPair {
	return updateRecurringTaskBulkOpPair{
		op: updateRecurringTaskBulkPairOp{
			Index: updateRecurringTaskBulkPairOpData{
				Id:            string(task.ID),
				Index:         TasquesRecurringTasksIndex,
				IfSeqNo:       uint64(task.Metadata.Version.SeqNum),
				IfPrimaryTerm: uint64(task.Metadata.Version.PrimaryTerm),
			},
		},
		doc: domainToPersistable(task, at),
	}
}

// <-- Persistence models

type jsonObjMap map[string]interface{}

type persistedRecurringTaskTaskDefinitionData struct {
	Queue             string        `json:"queue"`
	RetryTimes        uint          `json:"retry_times"`
	Kind              string        `json:"kind"`
	ProcessingTimeout time.Duration `json:"processing_timeout"`
	Priority          int           `json:"priority"`
	Args              *jsonObjMap   `json:"args,omitempty"`
	Context           *jsonObjMap   `json:"context,omitempty"`
}

type persistedRecurringTaskData struct {
	ScheduleExpression string                                   `json:"schedule_expression"`
	TaskDefinition     persistedRecurringTaskTaskDefinitionData `json:"task_definition"`
	IsDeleted          bool                                     `json:"is_deleted"`
	LoadedAt           *time.Time                               `json:"loaded_at,omitempty"`
	Metadata           common.PersistedMetadata                 `json:"metadata"`
}

func (e *EsService) newToPersistable(task *recurring.NewTask) persistedRecurringTaskData {
	now := e.getUTC()
	return persistedRecurringTaskData{
		ScheduleExpression: string(task.ScheduleExpression),
		TaskDefinition:     domainTaskDefToPersistable(&task.TaskDefinition),
		IsDeleted:          false,
		LoadedAt:           nil,
		Metadata: common.PersistedMetadata{
			CreatedAt:  now,
			ModifiedAt: now,
		},
	}
}

func domainToPersistable(task *recurring.Task, at metadata.ModifiedAt) persistedRecurringTaskData {
	return persistedRecurringTaskData{
		ScheduleExpression: string(task.ScheduleExpression),
		TaskDefinition:     domainTaskDefToPersistable(&task.TaskDefinition),
		IsDeleted:          bool(task.IsDeleted),
		LoadedAt:           (*time.Time)(task.LoadedAt),
		Metadata: common.PersistedMetadata{
			CreatedAt:  time.Time(task.Metadata.CreatedAt),
			ModifiedAt: time.Time(at),
		},
	}
}

func domainTaskDefToPersistable(def *recurring.TaskDefinition) persistedRecurringTaskTaskDefinitionData {
	return persistedRecurringTaskTaskDefinitionData{
		Queue:             string(def.Queue),
		RetryTimes:        uint(def.RetryTimes),
		Kind:              string(def.Kind),
		ProcessingTimeout: time.Duration(def.ProcessingTimeout),
		Priority:          int(def.Priority),
		Args:              (*jsonObjMap)(def.Args),
		Context:           (*jsonObjMap)(def.Context),
	}
}

func persistedToDomain(id task.RecurringTaskId, data *persistedRecurringTaskData, version metadata.Version) recurring.Task {
	return recurring.Task{
		ID:                 id,
		ScheduleExpression: recurring.ScheduleExpression(data.ScheduleExpression),
		TaskDefinition:     persistedTaskDefToDomainTaskDef(&data.TaskDefinition),
		IsDeleted:          recurring.IsDeleted(data.IsDeleted),
		LoadedAt:           (*recurring.LoadedAt)(data.LoadedAt),
		Metadata: metadata.Metadata{
			CreatedAt:  metadata.CreatedAt(data.Metadata.CreatedAt),
			ModifiedAt: metadata.ModifiedAt(data.Metadata.ModifiedAt),
			Version:    version,
		},
	}
}

func persistedTaskDefToDomainTaskDef(def *persistedRecurringTaskTaskDefinitionData) recurring.TaskDefinition {
	return recurring.TaskDefinition{
		Queue:             queue.Name(def.Queue),
		RetryTimes:        task.RetryTimes(def.RetryTimes),
		Kind:              task.Kind(def.Kind),
		Priority:          task.Priority(def.Priority),
		ProcessingTimeout: task.ProcessingTimeout(def.ProcessingTimeout),
		Args:              (*task.Args)(def.Args),
		Context:           (*task.Context)(def.Context),
	}
}

// persistence models -->

// <-- ES wrapped models

type esHitPersistedRecurringTask struct {
	ID          string                     `json:"_id"`
	Index       string                     `json:"_index"`
	SeqNum      uint64                     `json:"_seq_no"`
	PrimaryTerm uint64                     `json:"_primary_term"`
	Source      persistedRecurringTaskData `json:"_source"`
}

func (pTask *esHitPersistedRecurringTask) toDomainTask() recurring.Task {
	return persistedToDomain(task.RecurringTaskId(pTask.ID), &pTask.Source, pTask.Version())
}

func (pTask *esHitPersistedRecurringTask) Version() metadata.Version {
	return metadata.Version{
		SeqNum:      metadata.SeqNum(pTask.SeqNum),
		PrimaryTerm: metadata.PrimaryTerm(pTask.PrimaryTerm),
	}
}

type recurringTasksWithScrollId struct {
	ScrollId       string
	RecurringTasks []recurring.Task
}

type esSearchScrollingResponse struct {
	Hits struct {
		Hits []esHitPersistedRecurringTask `json:"hits"`
	} `json:"hits"`
	ScrollId string `json:"_scroll_id"`
}

// bulk

type updateRecurringTaskBulkOpPair struct {
	op  updateRecurringTaskBulkPairOp
	doc persistedRecurringTaskData
}

type updateRecurringTaskBulkPairOp struct {
	Index updateRecurringTaskBulkPairOpData `json:"index"`
}

type updateRecurringTaskBulkPairOpData struct {
	Id      string `json:"_id"`
	Index   string `json:"_index"`
	IfSeqNo uint64 `json:"if_seq_no"`

	IfPrimaryTerm uint64 `json:"if_primary_term"`
}

// 	Es wrapped models -->
