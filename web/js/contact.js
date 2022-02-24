const Endpoint = "https://gcodes.io/contact/submit";
const DevEndpoint = 'http://127.0.0.1:3000/contact/submit';


// submit payment success message to process order for fulfillment
function submitContactForm() {
  let formRaw = JSON.stringify($('#contact-form').serializeArray());
  let jsonArray = JSON.parse(formRaw);

  // validate input
  for (let i = 0; i < jsonArray.length; i++) {
      if (jsonArray[i].value == '') {
      console.log("NULL");
      return alert("Please make sure all information is filled out!")
      }
  }

  let payload = {
      email: jsonArray[0].value,
      subject:  jsonArray[1].value,
      message:  jsonArray[2].value,
  };

  // reset form
  document.querySelector('#contact-form').reset();

  // make db put call
  try {
    postRequest(Endpoint, payload)
    .then((response) => {
      return contactResponse(response.body)
    })
    .catch((err) => {
      console.log('err: ' + err)
      return errResponse(err);
    })
  } catch (error) {
    console.log('err: ' + error);
    return errResponse(error);
  }
}


// put call
async function postRequest(url = '', item = {}) {
  try {
    // Default options are marked with *
    const response = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
        body: JSON.stringify(item)
    });
    return response.json(); // parses JSON response into native JavaScript objects
  } catch (err) {
    console.log('postRequest err: ' + err);
    return err
  }
}


function contactResponse(resp) {
    return alert(resp.message)
}


function errResponse(err) {
    return alert(err)
}


function showModal(resp) {
  let modal = document.querySelector('#modal');
  let msg = document.querySelector('#error-msg');
  msg.innerHTML = resp.Message
  if (resp.Code == 'ITEM_OUT_OF_STOCK') {
    showFailed(resp.Failed);
  }
  modal.style.display = "block";
};


function showFailed(failed) {
  let div = document.querySelector('#failed');
  let ul = document.createElement('ul');
  ul.id = 'failed-list';
  for (let i = 0; i < failed.length; i++) {
    let li = document.createElement('li');
    li.innerHTML = failed[i].Name + ' (' + failed[i].Size + ')';
    ul.appendChild(li);
  }
  div.appendChild(ul);
}


function closeModal() {
  let modal = document.querySelector('#modal');
  let failed = modal.querySelector('#failed-list');
  failed.remove();
  
  modal.style.display = "none";
};
