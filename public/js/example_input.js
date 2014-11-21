
var examples = {
  healthcare: [
    'What is HIV?',
    'What are the benefits of taking aspirin daily?',
    'What can I do to get more calcium?',
    'How do I quit smoking?',
    'Who is at risk for diabetes?',
    'I am at risk for high blood pressure?'
  ],

  travel: [
    'Do I need a visa to enter Brazil?',
    'How much should I tip a taxi in Argentina?',
    'Where is the best place to dive in Austrlia?',
    'How high is Mount Everest?',
    'How deep is the Grand Canyon?',
    'When is the rainy season in Bangalore?'
  ]
};

function loadExample() {
  var corpus = $("#select").val();
  var index = Math.floor(Math.random() * examples[corpus].length);
  $('#questionText').val(examples[corpus][index]);
}

//fill and submit the form with a random example
function showExample(submit) {
  loadExample();
  if (submit)
    $('#qaForm').submit();
}

document.onload=($('#questionText').val() === '') ? loadExample() : '';