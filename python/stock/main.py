from datetime import datetime

import requests


def crawl_one_symbol(symbol, start_date=None, end_date=None):
    API_VNDIRECT = 'https://finfo-api.vndirect.com.vn/v4/stock_prices/'
    query = 'code:' + symbol
    params = {
        "sort": "date",
        "page": 1,
        "q": query
    }
    if start_date and end_date:
        query += '~date:gte:' + start_date + '~date:lte:' + end_date
        delta = datetime.strptime(end_date, '%Y-%m-%d') - \
            datetime.strptime(start_date, '%Y-%m-%d')
        params["size"]=delta.days + 1
    res = requests.get(API_VNDIRECT, params=params)
    res.raise_for_status()
    data = res.json()['data']
    return data


crawl_one_symbol("TCB")
