const container = document.querySelector('.containerTable');


const renderTable = function(data) {   
    let userCountView = 0;
    data.forEach(user => {
        userCountView++;
        const lastUser = document.querySelector('tr:last-of-type');
        const markup = `
            <tr>
                <td>${userCountView}</td>
                <td>${user.RoleID}</td>
                <td>${user.Name}</td>
                <td>${user.Surname}</td>
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