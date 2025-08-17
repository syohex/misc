namespace BlogGen;

public static class ParserFactory
{
    public static IParser Create(string url)
    {
        if (url.Contains("dmm.co.jp"))
        {
            return new DmmParser();
        }

        throw new Exception($"Unsupported site: {url}");
    }
}