using TCBSClient;

var client = new APIClient();
var company = await client.GetCompanyAsync("TCB");
Console.WriteLine(company.ShortName);
