function search(searchDB, searchResultModal) {
    var keyword = $("#search").val();
    if (keyword.length >= 2) {
        var searchResult = searchDB.search(keyword);
        if (searchResult.length > 0) {
            var searchResultBox = $("#search-result-box");
            searchResultBox.html("");
            searchResult.forEach(function(val, index, arr) {
                if (val.doc.title == "About") {
                    return
                }
                var top_margin = 2;
                if (index == 0) {
                    top_margin = 0;
                }
                searchResultBox.append(
                    '<div class="mt-' + top_margin + '">' +
                    '<a href="' + val.doc.id + '">' +
                    '<div class="article-title m-0 p-0 d-inline-flex">' +
                    '<h1 class="article-title-text article-title-radius px-2 py-1 m-0">' +
                    val.doc.title +
                    '</h1>' +
                    '</div>' +
                    '</a>' +
                    '<div class="article-body">' +
                    val.doc.body.substring(1, 150) + '...' +
                    '</div>' +
                    '</div>'
                );
            });
        } else {
            $("#search-result-box").html("NOT FOUND!");
        }
        searchResultModal.toggle();
    }
}

$(document).ready(function() {
    // Initialize the search by indexing data.
    var searchDB = elasticlunr.Index.load(window.searchIndex);
    const searchResultModal = new bootstrap.Modal("#search-result");

    // by clicking on search button.
    $("#search-btn").click(function() {
        search(searchDB, searchResultModal);
    });

    // by pressing Enter on the search box.
    $("#search").keypress(function(e) {
        if (e.which == 13) {
            search(searchDB, searchResultModal);
        }
    });
})
