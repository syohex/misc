using Microsoft.Playwright;

class Program
{
    public static async Task Main()
    {
        const string url = "https://video.dmm.co.jp/av/content/?id=sone00638";
        using var playwright = await Playwright.CreateAsync();
        var browser = await playwright.Chromium.LaunchAsync(new BrowserTypeLaunchOptions { Headless = true });
        var context = await browser.NewContextAsync();

        var cookie = new Cookie
        {
            Name = "age_check_done",
            Value = "1",
            Path = "/",
            Domain = ".dmm.co.jp",
            HttpOnly = true,
            Secure = true,
        };
        await context.AddCookiesAsync([cookie]);

        var page = await context.NewPageAsync();
        await page.GotoAsync(url);

        var title = await page.TitleAsync();
        Console.WriteLine($"Page Title: {title}");

        var handles = await page.Locator("a").ElementHandlesAsync();
        foreach (var handle in handles)
        {
            var src = await handle.GetAttributeAsync("href");
            if (src == null) continue;

            if (src.EndsWith("pl.jpg"))
            {
                Console.WriteLine($"Found Image: {src}");
                break;
            }
        }

        await context.CloseAsync();
        await browser.CloseAsync();
    }
}