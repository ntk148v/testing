namespace TCBSClient;

public readonly struct Company
{
    public string Exchange { get; }
    public string ShortName { get; }
    public int IndustryID { get; }
    public string TndustryIDv2 { get; }
    public string Tndustry { get; }
    public string TndustryEn { get; }
    public string EstablishedYear { get; }
    public int NoEmployees { get; }
    public int NoShareholders { get; }
    public float ForeignPercent { get; }
    public Uri Website { get; }
    public float StockRating { get; }
    public float DeltaInWeek { get; }
    public float DeltaInMonth { get; }
    public float DeltaInYear { get; }
    public float OutstandingShare { get; }
    public float IssueShare { get; }
    public string CompanyType { get; }
    public string Ticker { get; }
}