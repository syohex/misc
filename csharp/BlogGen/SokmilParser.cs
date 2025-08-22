using Microsoft.AspNetCore.WebUtilities;
using Microsoft.Playwright;

namespace BlogGen;

public class SokmilParser : IParser
{
    public async Task<Product> Parse(string url, Config config)
    {
        using var playwright = await Playwright.CreateAsync();
        var browser = await playwright.Chromium.LaunchAsync(new BrowserTypeLaunchOptions { Headless = false });
        var context = await browser.NewContextAsync();

        var cookie = new Cookie
        {
            Name = "AGEAUTH",
            Value = "ok",
            Path = "/",
            Domain = ".sokmil.com",
            HttpOnly = true,
            Secure = true,
        };
        await context.AddCookiesAsync([cookie]);

        var page = await context.NewPageAsync();
        await page.GotoAsync(url);

        var product = new Product();

        var titleElement = await page.QuerySelectorAsync("h1.page-title");
        if (titleElement != null)
        {
            var title = await titleElement.TextContentAsync();
            if (title == null)
            {
                throw new Exception($"Cannot get title from {url}");
            }

            product.Title = title;
        }

        var packageElement = await page.QuerySelectorAsync("a.sokmil-lightbox-jacket");
        if (packageElement != null)
        {
            var title = await packageElement.GetAttributeAsync("href");
            if (title == null)
            {
                throw new Exception($"Cannot get package from {url}");
            }

            product.Image = title;
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
            ["affi"] = config.Sokmil.Id,
            ["utm_source"] = "sokmil_ad",
            ["utm_medium"] = "affiliate",
            ["utm_campaign"] = config.Sokmil.Id
        };

        var builder = new UriBuilder(url)
        {
            Query = string.Empty,
            Fragment = string.Empty
        };

        return QueryHelpers.AddQueryString(builder.Uri.ToString(), queries);
    }
}