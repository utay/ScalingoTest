var COLORS = ['#0074D9', '#7FDBFF', '#39CCCC', '#3D9970', '#2ECC40', '#01FF70', '#FFDC00'];
var $repositories = $('#repositories');
var template =
    '<tr>' +
    '<td>{{ID}}</td>' +
    '<td><a href="https://github.com/{{Name}}">{{Name}}</a></td>' +
    '<td><ul>{{#Languages}}' +
    '<li>{{Name}}: {{Lines}}</li>' +
    '{{/Languages}}</ul></td>' +
    '</tr>';

Mustache.parse(template);

var byLanguage = {};
$.each(window.repositories, function(i, repo) {
    $.each(repo.Languages, function(j, lang) {
        byLanguage[lang.Name] = byLanguage[lang.Name] || [];
        byLanguage[lang.Name].push(repo);
    });
});
console.log(byLanguage);

function renderTable() {
  var html = '';
  $.each(window.repositories, function(i, repo) {
    html += Mustache.render(template, repo);
  });
  $('#repositories').html(html);
}

function renderChartByLanguage() {
  var data = [];
  var i = 0;
  for (var lang in byLanguage) {
    data.push({
      value: byLanguage[lang].length,
      label: lang,
      color: COLORS[i++ % COLORS.length]
    });
  }
  var ctx = document.getElementById("by-language").getContext("2d");
  new Chart(ctx).Pie(data);
}

function renderChartLinesPerLanguage() {
    var linesByLanguage = {};
    $.each(window.repositories, function(i, repo) {
        $.each(repo.Languages, function(j, lang) {
            linesByLanguage[lang.Name] = linesByLanguage[lang.Name] || 0;
            linesByLanguage[lang.Name] += lang.Lines;
        });
    });

    var data = [];
    var i = 0;
    for (var name in linesByLanguage) {
        data.push({
            value: linesByLanguage[name],
            label: name,
            color: COLORS[i++ % COLORS.length]
        });
    }
  var ctx = document.getElementById("lines-per-language").getContext("2d");
  new Chart(ctx).Pie(data);
}

$('.sort-by').on('click', function(event) {
  event.preventDefault();
  var attribute = $(this).data('attr');
  var order = $(this).data('order');
  var number = $(this).data('number')

  window.repositories.sort(function(a, b) {
    var first = order == 'desc' ? a[attribute] : b[attribute];
    var second = order == 'desc' ? b[attribute] : a[attribute];
    if (number) {
      first = +first;   // cast to int
      second = +second; // cast to int
      return first < second ? -1 : (first > second ? 1 : 0);
    }
    // string
    first = first || '';
    second = second || '';
    return first.localeCompare(second);
  });
  renderTable();
});

renderTable();
renderChartByLanguage();
renderChartLinesPerLanguage();

