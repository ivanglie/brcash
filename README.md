# brcash

A service that displays up-to-date cash currency exchange rates in Russian banks.

## Example

Output:

```json
{
  "currency": "USD",
  "city": "novosibirsk",
  "items": [
    {
      "bank": "ДО \"Новосибирский\" Ф-ла Сибирский ПАО Банк \"ФК Открытие\"",
      "subway": "Октябрьская, Речной вокзал, Площадь Ленина",
      "buy": 87.5,
      "sell": 90.8,
      "updated": "2023-07-16T18:07:00+03:00"
    },
    {
      "bank": "ДО \"Новосибирск\"",
      "subway": "Площадь Ленина, Красный проспект, Площадь Гарина-Михайловского",
      "buy": 87,
      "sell": 90.6,
      "updated": "2023-07-16T18:07:00+03:00"
    },
    {
      "bank": "\"АК БАРС\" Банк - пр. Димитрова, 7",
      "subway": "Площадь Ленина, Красный проспект, Площадь Гарина-Михайловского",
      "buy": 86.89,
      "sell": 90.91,
      "updated": "2023-07-16T18:06:00+03:00"
    }
  ]
}
```

## References

For more information check out the following links:

* Cash currency exchange rates by [Banki.ru](https://www.banki.ru/products/currency/map/moskva/) (RU)
