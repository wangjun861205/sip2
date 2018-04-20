# sip2
## Overview
Restful style communicate server for 3M SIPII protocal

## JSON request format
```{
      "header": {
        "method": "query_patron_status"
      },
      "data": {
        "language": 1,
        "transaction_date": "2018-04-16 15:07:01",
        "institution_id": "0001",
        "patron_id": "0001",
        "terminal_password": "terminal password",
        "patron_password": "patron password",
       }
   }
```

## JSON response format
```
  {
    "header": {
      "version": "1.00",
      },
    "data": {
      "msg": "ok",
      "code": 200,
      "item_list": null,
      "item": {
         "patron_status": "xxxxx",
         "language": 1,
         ...
         }
      "meta": null
      }
   }
```
      
