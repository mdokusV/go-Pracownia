const container = document.querySelector('.containerTable');


const renderTable = function(data) {   
    let userCountView = 0;
    data.forEach(user => {
        userCountView++;
        const lastUser = document.querySelector('tr:last-of-type');
        const markup = `
            <tr>
                <td>${userCountView}</td>
                <td>${user.ID}</td>
                <td>${user.RoleID}</td>
                <td>${user.Name}</td>
                <td>${user.Surname}</td>
                <td>${user.DateOfBirth}</td>
                <td>${user.Login}</td>
                <td>${user.Password}</td>
                <td>${user.CreatedAt}</td>
                <td>${user.UpdatedAt}</td>
                <td>${user.DeletedAt}</td>
                <td><button class="btn btn-primary btn-removeUser">Remove</button></td>
            </tr>
        `;
        lastUser.insertAdjacentHTML('afterend', markup);
    });
}

window.addEventListener('load', async () => {
    try {
        const result = await fetch('/UserShowAll');
        const data = await result.json();
        renderTable(data)
    } catch(e) {
        console.error(e)
    }
})