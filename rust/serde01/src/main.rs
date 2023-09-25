use std::fs;
use serde::Deserialize;

#[derive(Deserialize, Debug)]
struct Data {
    name: String,
    age: i32,
    description: String,
}

fn main() -> anyhow::Result<()>{
    let data = fs::read_to_string("test.json")?;
    let d: Vec<Data> = serde_json::from_str(&data)?;
    println!("d={d:?}");

    Ok(())
}
