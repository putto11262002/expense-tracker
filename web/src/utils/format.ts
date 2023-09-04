export function formatAsCurrency(number: number, currencyCode = "USD") {
    return number.toLocaleString("en-US", {
      style: "currency",
      currency: currencyCode,
    });
  }