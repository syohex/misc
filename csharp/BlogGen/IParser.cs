namespace BlogGen;

public interface IParser
{
    Task<Product> Parse(string url, Config config);
}