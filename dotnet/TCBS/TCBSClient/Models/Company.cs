namespace TCBSClient;

public struct Company
{
    public string Exchange { get; set; }
    public string ShortName { get; set; }
    public int IndustryID { get; set; }
    public string TndustryIDv2 { get; set; }
    public string Tndustry { get; set; }
    public string TndustryEn { get; set; }
    public string EstablishedYear { get; set; }
    public int NoEmployees { get; set; }
    public int NoShareholders { get; set; }
    public float ForeignPercent { get; set; }
    public Uri Website { get; set; }
    public float StockRating { get; set; }
    public float DeltaInWeek { get; set; }
    public float DeltaInMonth { get; set; }
    public float DeltaInYear { get; set; }
    public float OutstandingShare { get; set; }
    public float IssueShare { get; set; }
    public string CompanyType { get; set; }
    public string Ticker { get; set; }
}