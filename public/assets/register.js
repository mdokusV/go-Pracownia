const form = document.querySelector('form');

form.addEventListener('submit', async (event) => {
  event.preventDefault();
 
  try {
    console.log(form.elements);
    const response = await fetch('/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        Login: form.elements.login.value,
        Password: form.elements.password.value,
        DateOfBirth: form.elements.dateofbirth.value,
        Name: form.elements.name.value,
        Surname: form.elements.surname.value,
      })
    });

    if (response.redirected) {
      // Server redirected you to another page and everything is fine
      window.location.href = response.url;
    } else {
      // Server sent you back JSON with variable Error in it
      const data = await response.json();
    
      alert('Something went wrong!');

      if (Array.isArray(data)) {
        data.forEach(err => console.error(err))
      } else {
        console.error(data)
      }

    }
  } catch (error) {
    console.error('Request failed', error);
  }
});