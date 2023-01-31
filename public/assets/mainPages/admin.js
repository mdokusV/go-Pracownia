const container = document.querySelector('.containerTable');
const btnNext = document.querySelector('#btn--next');
const btnPrevious = document.querySelector('#btn--previous');

const state = {
    pageNumber: 1,
    maxPages: 1,
}

const clearTable = function() {
    container.innerHTML = "";
    renderInitialTable();
}

const renderInitialTable = function() {
    const markup = `
    <table class="table">
        <tr id="table-headers">
            <th>LP</th>
            <th>ID</th>
            <th>RoleName</th>
            <th>Change Role</th>
            <th>Name</th>
            <th>Surname</th>
            <th>Date of Birth</th>
            <th>Login</th>
            <th>Password</th>
            <th>CreatedAt</th>
            <th>UpdatedAt</th>
            <th>DeletedAt</th>
            <th>Delete:</th>
        </tr>
    </table>
    `;
    container.insertAdjacentHTML('afterbegin', markup);
}

const queryPages = async function() {
    try {
        // Query page
        const pageNumberJSON = JSON.stringify({
            pageNumber: state.pageNumber,
        });

        // Query maxPage
        const resMaxPage = await fetch('/MaxPages');
        const resPageNum = await fetch('/UserShow', {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
            body: pageNumberJSON,
        });
        const dataMaxPage = await resMaxPage.json();
        const dataPageNum = await resPageNum.json();
        
        renderTable(dataPageNum)
        state.maxPages = dataMaxPage.MaxPages;

    } catch(e) {
        console.error(e)
    }
}

btnNext.addEventListener('click', async () => {
        if(state.pageNumber > state.maxPages-1) return;
        state.pageNumber++;
        await queryPages();
})

btnPrevious.addEventListener('click', async () => {
    if(state.pageNumber-1 < 1) return;
    state.pageNumber--;
    await queryPages();
})



const renderTable = function(data) {
    clearTable();
    data.forEach(user => {
        const lastUser = document.querySelector('tr:last-of-type');
        const markup = `
            <tr>
                <td>${user.OrderNumber}</td>
                <td>${user.ID}</td>
                <td>${user.RoleName}</td>
                <td><select name="roleName" data-login="${user.Login}" class="form-select">
                    <option value="user" ${user.RoleName === 'user' && 'selected'}>user</option>
                    <option value="moderator" ${user.RoleName === 'moderator' && 'selected'}>moderator</option>
                    <option value="admin" ${user.RoleName === 'admin' && 'selected'}>admin</option>
                </select></td>
                <td>${user.Name}</td>
                <td>${user.Surname}</td>
                <td>${user.DateOfBirth}</td>
                <td>${user.Login}</td>
                <td>${user.Password}</td>
                <td>${user.CreatedAt}</td>
                <td>${user.UpdatedAt}</td>
                <td>${user.DeletedAt}</td>
                <td><button data-login="${user.Login}" class="btn btn-primary btn-removeUser">Remove</button></td>
            </tr>
        `;
        lastUser.insertAdjacentHTML('afterend', markup);
    });
}

container.addEventListener('click', async e => {
    if(e.target.tagName !== 'BUTTON') return;
    
    try {
        const closestBtn = e.target.closest('.btn-removeUser').getAttribute('data-login');
        const dataForSendJSON = JSON.stringify({
            login: `${closestBtn}`,
        });
        
        const result = await fetch('/UserDelete', {
            method: "DELETE",
            headers: {
                'Content-Type': 'application/json',
            },
            body: dataForSendJSON,
        });

        await queryPages();

        
    } catch (e) {
        console.error(e);
    }

})

container.addEventListener('change', async e => {
    if(e.target.tagName !== 'SELECT') return;
    try {
        const login = e.target.closest('select').getAttribute('data-login');
        const dataForSendJSON = JSON.stringify({
            Login: `${login}`,
            RoleName: `${e.target.closest('select').value}`,
        });


        
        await fetch('/changeUserRole', {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
            },
            body: dataForSendJSON,
        });
        await queryPages();

        
    } catch (e) {
        console.error(e);
        await queryPages();
    }
})
window.addEventListener('load', queryPages);




