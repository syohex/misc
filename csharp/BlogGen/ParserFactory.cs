namespace BlogGen;

public static class ParserFactory
{
    public static IParser Create(string url)
    {
        if (url.Contains("dmm.co.jp"))
        {
            return new DmmParser();
        }

        if (url.Contains("sokmil.com"))
        {
            return new SokmilParser();
        }

        throw new Exception($"Unsupported site: {url}");
    }
}