package testdata

const TransactionWithResults = `
{
    "ref_block_num": 56327,
    "ref_block_prefix": 3761836804,
    "expiration": "2018-09-17T08:42:05",
    "operations": [
      [
        2,
        {
          "fee": {
            "amount": 0,
            "asset_id": "1.3.0"
          },
          "fee_paying_account": "1.2.19",
          "order": "1.7.4442286",
          "extensions": []
        }
      ],
      [
        2,
        {
          "fee": {
            "amount": 0,
            "asset_id": "1.3.0"
          },
          "fee_paying_account": "1.2.19",
          "order": "1.7.4442287",
          "extensions": []
        }
      ],
      [
        2,
        {
          "fee": {
            "amount": 0,
            "asset_id": "1.3.0"
          },
          "fee_paying_account": "1.2.19",
          "order": "1.7.4442288",
          "extensions": []
        }
      ],
      [
        1,
        {
          "fee": {
            "amount": 10000,
            "asset_id": "1.3.0"
          },
          "seller": "1.2.19",
          "amount_to_sell": {
            "amount": "100000000000",
            "asset_id": "1.3.0"
          },
          "min_to_receive": {
            "amount": "10535556799",
            "asset_id": "1.3.8"
          },
          "expiration": "2018-09-17T08:43:37",
          "fill_or_kill": false,
          "extensions": []
        }
      ],
      [
        1,
        {
          "fee": {
            "amount": 10000,
            "asset_id": "1.3.0"
          },
          "seller": "1.2.19",
          "amount_to_sell": {
            "amount": "100000000000",
            "asset_id": "1.3.0"
          },
          "min_to_receive": {
            "amount": "10638342719",
            "asset_id": "1.3.8"
          },
          "expiration": "2018-09-17T08:43:37",
          "fill_or_kill": false,
          "extensions": []
        }
      ],
      [
        1,
        {
          "fee": {
            "amount": 10000,
            "asset_id": "1.3.0"
          },
          "seller": "1.2.19",
          "amount_to_sell": {
            "amount": "100000000000",
            "asset_id": "1.3.0"
          },
          "min_to_receive": {
            "amount": "10741128639",
            "asset_id": "1.3.8"
          },
          "expiration": "2018-09-17T08:43:37",
          "fill_or_kill": false,
          "extensions": []
        }
      ]
    ],
    "extensions": [],
    "signatures": [
      "20594600f763d8f74e4d7dee0f16867c4431d091b04712dfe2fe75f9a7747209683f6ad2e29a95fa9842b45dafc3159db54f27f206ae4792e0cb4d23dd175f0efd"
    ],
    "operation_results": [
      [
        2,
        {
          "amount": "100000000000",
          "asset_id": "1.3.0"
        }
      ],
      [
        2,
        {
          "amount": "100000000000",
          "asset_id": "1.3.0"
        }
      ],
      [
        2,
        {
          "amount": "100000000000",
          "asset_id": "1.3.0"
        }
      ],
      [
        1,
        "1.7.4442340"
      ],
      [
        1,
        "1.7.4442341"
      ],
      [
        1,
        "1.7.4442342"
      ]
    ]
  }
`
