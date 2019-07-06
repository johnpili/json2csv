# json2csv

## Overview
I got multiple projects where the customers needed CSV rather than JSON. My existing my microservice is in RESTful and to accommodate the CSV requirement I created this Golang based JSON to CSV tool. Because this tool is intended for more than one customer and will have to use different endpoints with different CSV output I have added a rule file or template mechanism to handle output formatting.

```json
[
    {
        "ruleName": "rule-name1",
        "ioMapping": [
            {
                "incoming": "dateTime",
                "outgoing": "DATETIME"
            },            
            {
                "incoming": "equipmentName",
                "outgoing": "EQUIPMENT"
            },
            {
                "incoming": "current",
                "outgoing": "CURRENT"
            },
            {
                "incoming": "power", 
                "outgoing": "POWER"
            },
            {
                "incoming": "reactivePower", 
                "outgoing": "REACTIVE_POWER"
            },
            {
                "incoming": "voltageReading", 
                "outgoing": "VOLTAGE"
            },
            {
                "incoming": "voltageType",
                "outgoing": "VOLTAGE_LEVEL"
            },
            {
                "incoming": "utilisation",
                "outgoing": "UTILISATION"
            },
            {
                "incoming": "substationType",
                "outgoing": "SS_TYPE"
            },
            {
                "incoming": "feederType",
                "outgoing": "FDR"
            }
        ],
        "targetUrlToken": ""
    }    
]
```