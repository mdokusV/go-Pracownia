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
            <th>RoleName</th>
            <th>Name</th>
            <th>Surname</th>
            <th>Date of Birth</th>
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
                <td>${user.RoleName}</td>
                <td>${user.Name}</td>
                <td>${user.Surname}</td>
                <td>${user.DateOfBirth}</td>
            </tr>
        `;
        lastUser.insertAdjacentHTML('afterend', markup);
    });
}



window.addEventListener('load', queryPages);


