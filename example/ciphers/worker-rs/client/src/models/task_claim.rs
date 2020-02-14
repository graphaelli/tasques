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
pub struct TaskClaim {
    #[serde(rename = "amount", skip_serializing_if = "Option::is_none")]
    pub amount: Option<i32>,
    #[serde(rename = "block_for", skip_serializing_if = "Option::is_none")]
    pub block_for: Option<String>,
    #[serde(rename = "queues")]
    pub queues: Vec<String>,
}

impl TaskClaim {
    pub fn new(queues: Vec<String>) -> TaskClaim {
        TaskClaim {
            amount: None,
            block_for: None,
            queues,
        }
    }
}