function addCake(id, name, price, amount) {
  const cake = document.querySelector(`[name="cake[${id}]"]`);
  if (cake != null) {
    cake.value = Number(cake.value) + 1;
    let tr = cake.parentElement.parentElement.parentElement;
    //console.log(tr);
    tr.querySelector(".cake-total").innerText =
      `${price * Number(cake.value)}PLN`;
    updateTotalPrice();
    return;
  }
  let template = document.getElementById("basket_element_template");
  let tr = template.content.querySelector("tr").cloneNode(true);
  tr.querySelector(".cake-name").innerText = name;
  tr.querySelector(".cake-price").innerText = `${price}PLN`;
  tr.querySelector(".cake-total").innerText = `${price * amount}PLN`;
  tr.querySelector("button").onclick = () => {
    tr.remove();
  };
  let input = tr.querySelector("input");
  input.onchange = () => {
    if (input.value == 0) {
      tr.remove();
    }
    tr.querySelector(".cake-total").innerText = `${price * input.value}PLN`;
    updateTotalPrice();
  };
  input.value = amount;
  input.name = `cake[${id}]`;
  document.querySelector("#basket_table").appendChild(tr);
  updateTotalPrice();
}
function updateTotalPrice() {
  document.querySelector("#total-price").innerText = `${totalPrice()}PLN`;
}
function totalPrice() {
  let out = 0;
  const cakes = document.querySelectorAll("#basket_table>tr");
  cakes.forEach((e) => {
    let priceString = e.querySelector(".cake-price").innerText;
    let price = Number(priceString.substring(0, priceString.length - 3));
    let amount = Number(e.querySelector("input").value);
    out += price * amount;
    // console.log(out);
  });
  return out;
}
function toggleDetails(row) {
  const detailsRow = row.nextElementSibling;
  if (detailsRow.classList.contains("hidden")) {
    detailsRow.classList.remove("hidden");
    return;
  }
  detailsRow.classList.add("hidden");
}
