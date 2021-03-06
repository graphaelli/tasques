/*
 * Tasques API
 *
 * A Task queue backed by Elasticsearch
 *
 * The version of the OpenAPI document: 0.0.1
 *
 * Generated by: https://openapi-generator.tech
 */

#[derive(Clone, Debug, PartialEq, Serialize, Deserialize)]
pub struct TaskReport {
    #[serde(rename = "at")]
    pub at: String,
    #[serde(rename = "data", skip_serializing_if = "Option::is_none")]
    pub data: Option<serde_json::Value>,
}

impl TaskReport {
    pub fn new(at: String) -> TaskReport {
        TaskReport { at, data: None }
    }
}
