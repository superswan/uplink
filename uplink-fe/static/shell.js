document.getElementById('submit').addEventListener('click', function () {
    var input = document.getElementById('input');
    var output = document.getElementById('output');

    // Print the command
    output.textContent += '> ' + input.value + '\n';

    // Clear the input
    input.value = '';

    // Scroll to the bottom of the output
    output.parentElement.scrollTop = output.parentElement.scrollHeight;
});