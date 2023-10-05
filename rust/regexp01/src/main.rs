use std::borrow::Cow;
use arboard::Clipboard;
use regex::Regex;

fn main() -> anyhow::Result<()> {
    let args: Vec<String> = std::env::args().collect();
    if args.len() < 5 {
        println!("Usage: {} pattern series_id series_no count", args[0]);
        return Ok(());
    }

    let count: usize = args[4].parse()?;

    let (pattern, series_id, series_no) = (&args[1], &args[2], &args[3]);
    let id_pattern = format!("(?P<name>(?i:{})){}", series_id, series_no);
    let re = Regex::new(&id_pattern)?;

    let initial_num: usize = series_no.parse()?;

    let mut strs: Vec<String> = vec![];
    for i in 1..=count {
        let num = initial_num + i;
        let replaced = format!("${{name}}{num}");
        let ret = re.replace_all(&pattern, replaced);
        println!("{ret}");

        match ret {
            Cow::Borrowed(s) => {
                strs.push(s.to_string());
            }
            _ => {}
        }
    }

    let output = strs.join("\n");
    let mut clipboard = Clipboard::new()?;
    clipboard.set_text(&output)?;

    Ok(())
}
