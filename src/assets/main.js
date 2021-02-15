var subtotal = 0;

function updateSubtotal() {
  var total = document.getElementsByClassName('total-price');
  var subtotal = 0;
  
  for (var i=0; i<total.length; i++) {
    subtotal = subtotal + parseFloat(total[i].dataset.total);
  }
  document.getElementById('subtotal-amount').innerHTML = subtotal.toFixed(2).replace(/\d(?=(\d{3})+\.)/g, '$&,');
  document.getElementById('edt-subtotal-amount-disp').value = subtotal.toFixed(2).replace(/\d(?=(\d{3})+\.)/g, '$&,');
  document.getElementById('edt-subtotal-amount').value = subtotal;
}

function countTotalPrice(e = null) {
  var qty = document.querySelectorAll('.edt-qty');
  for (var i=0; i<qty.length; i++) {          
    var dataId = qty[i].dataset.id;
    var dataQty = qty[i].value;
    var dataStock = document.getElementById('stock-'+dataId).dataset.stock;
    var dataPrice = document.getElementById('price-'+dataId).dataset.price;
    var totalPrice = dataQty * dataPrice;

    if (parseInt(dataQty) <= parseInt(dataStock)) {
      document.getElementById('total-'+dataId).value = totalPrice.toFixed(2).replace(/\d(?=(\d{3})+\.)/g, '$&,');
      document.getElementById('total-'+dataId).dataset.total = totalPrice;                 
      document.getElementById('skuqty-'+dataId).value = dataId+'-'+dataQty;
      qty[i].style.borderColor = "#888";
    } else {
      qty[i].style.borderColor = "red";
    }
  }
  updateSubtotal();
};

function processCheckout(elm) {
  elm.disabled = true;
  document.getElementById('frmCheckout').submit();
}

window.onload = function() {            
  countTotalPrice();
}