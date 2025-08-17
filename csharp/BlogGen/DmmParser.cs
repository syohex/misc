using Microsoft.AspNetCore.WebUtilities;
using Microsoft.Playwright;

namespace BlogGen;

public class DmmParser : IParser
{
    private const string AffiliateUrl = "https://al.dmm.co.jp";

    public async Task<Product> Parse(string url, Config config)
    {
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

        var product = new Product();

        var metaTitle = await page.QuerySelectorAsync("meta[property='og:title']");
        if (metaTitle != null)
        {
            var title = await metaTitle.GetAttributeAsync("content");
            if (title == null)
            {
                throw new Exception($"Cannot get title from {url}");
            }

            product.Title = title;
        }

        var handles = await page.Locator("a").ElementHandlesAsync();
        foreach (var handle in handles)
        {
            var src = await handle.GetAttributeAsync("href");
            if (src == null) continue;

            if (src.EndsWith("pl.jpg"))
            {
                product.Image = src;
                break;
            }
        }

        await context.CloseAsync();
        await browser.CloseAsync();

        product.Url = ConstructAffiliateUrl(url, config);
        return product;
    }

    private string ConstructAffiliateUrl(string url, Config config)
    {
        var queries = new Dictionary<string, string?>
        {
            ["lurl"] = url,
            ["af_id"] = config.Dmm.Id,
            ["ch"] = "link_tool",
            ["ch_id"] = "link",
        };

        return QueryHelpers.AddQueryString(AffiliateUrl, queries);
    }
}