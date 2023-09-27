fn main() {
    let response = reqwest::blocking::get("https://syohex.org").unwrap();
    let content = response.text().unwrap();

    let doc = scraper::Html::parse_document(&content);
    let selector = scraper::Selector::parse("a").unwrap();
    let links = doc.select(&selector);

    for link in links {
        let attr = link.value().attr("href").unwrap();
        println!("link={attr}");
    }
}
