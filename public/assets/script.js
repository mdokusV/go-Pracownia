const btnTable = document.querySelector('#button--table');
const btnLogin = document.querySelector('#button--login');

const container = document.querySelector('.container');

const state = {
    posts: "",
}


btnTable.addEventListener('click', async () => {
    try {
        const result = await fetch('/ShowAllPosts');
        const data = await result.json();
        state.posts = data;

        renderTable(state.posts)
    } catch(e) {
        console.error(e)
    }
})




const renderTable = function(data) {
    clearContainer()

    const table = `
        <table class="table">
            <tr>
            <th>ID</th>
            <th>Post Title</th>
            <th>Post Body</th>
            </tr>
        </table>
    `;
    
    container.insertAdjacentHTML('afterbegin', table);
    

    data.message.forEach(post => {
        const lastPost = document.querySelector('tr:last-of-type');
        const markup = `
            <tr>
                <td>${post.ID}</td>
                <td>${post.Title}</td>
                <td>${post.Body}</td>
            </tr>
        `;

        lastPost.insertAdjacentHTML('afterend', markup);
    });
}

const clearContainer = function() {
    container.innerHTML = "";
}

const renderLogin = function() {
    clearContainer();
    const form = `
        <form class="form" action="#" method="POST">
            <label for="form__login">Login</label>
            <input id="form__login" type="email" placeholder="Enter login" required>
            <label for="form__password">Password</label>
            <input id="form__password" type="password" placeholder="Enter password" required>
            <input type="submit" value="Login">
        </form>
    `;
    container.insertAdjacentHTML('afterbegin', form);
}



btnLogin.addEventListener('click', renderLogin);