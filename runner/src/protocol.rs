#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Entity {
    pub id: u64,
    pub name: Option<String>,
    pub color: String,
    pub sprite: Option<String>,
    pub x: f64,
    pub y: f64,
    pub radius: f64,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct WorldState {
    pub time: u64,
    pub entities: Vec<Entity>,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct TickRequest {
    pub world_state: WorldState,
    pub player: Entity,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct TickResult {
    pub dx: f64,
    pub dy: f64,
}
