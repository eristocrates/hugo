var data = {};
var audio = {};
var hits_correct = 0;
var hits_wrong = 0;
var start_time = 0;
var hpm = 0;
var ratio = 0;
var isModalOpen = false;


data.chars = " jfkdlsahgyturieowpqbnvmcxz6758493021`-=[]\\;',./ABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+{}|:\"<>?";
data.consecutive = 5;
data.word_length = 7;
data.current_layout = "qwerty";
const defaultLayouts = {
    "qwerty": " jfkdlsahgyturieowpqbnvmcxz6758493021`-=[]\\;',./ABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+{}|:\"<>?",
    "azerty": " jfkdlsmqhgyturieozpabnvcxw6758493021`-=[]\\;',./ABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+{}|:\"<>?",
    "colemak": " ntesiroahdjglpufywqbkvmcxz1234567890'\",.!?:;/@$%&#*()_ABCDEFGHIJKLMNOPQRSTUVWXYZ~+-={}|^<>`[]\\",
    "bepo": " tesirunamc,èvodpléjbk'.qxghyfàzw6758493021`-=[]\\;/ABCDEFGHIJKLMNOPQRSTUVWXYZ~!@#$%^&*()_+{}|:\"<>?",
    "norman": " ntieosaygjkufrdlw;qbpvmcxz1234567890'\",.!?:;/@$%&#*()_ABCDEFGHIJKLMNOPQRSTUVWXYZ~+-={}|^<>`[]\\",
    "code-es6": " {}',;():.>=</_-|`!?#[]\\+\"@$%&*~^"
};
layouts = { ...defaultLayouts };
// bépo
function stringDivide(input) {
    const length = input.length;
    const midpoint = Math.floor(length / 2);

    let left = midpoint - 1;
    let right = midpoint;
    if (length % 2 === 0) {
        right = midpoint;
    } else {
        right = midpoint + 1;
    }

    let result = '';
    while (left >= 0 || right < length) {
        if (right < length) {
            result += input[right];
            right++;
        }
        if (left >= 0) {
            result += input[left];
            left--;
        }
    }

    return result;
}
function load_layouts() {
    const savedLayouts = localStorage.getItem('layouts');
    if (savedLayouts) {
        layouts = JSON.parse(savedLayouts);
    }
}

function save_layouts() {
    localStorage.setItem('layouts', JSON.stringify(layouts));
}
function delete_current_layout() {

    const layoutKeys = Object.keys(layouts);

    // Prevent deletion if there's only one layout left
    if (layoutKeys.length <= 1) {
        alert("Cannot delete the last remaining layout.");
        return;
    }



    if (data.current_layout && layouts[data.current_layout]) {
        const layoutKeys = Object.keys(layouts);
        const currentIndex = layoutKeys.indexOf(data.current_layout);

        // Delete the current layout
        delete layouts[data.current_layout];

        // Determine the new current layout
        let newLayout;
        if (currentIndex > 0) {
            newLayout = layoutKeys[currentIndex - 1];
        } else if (layoutKeys.length > 1) {
            newLayout = layoutKeys[1];
        } else {
            newLayout = null;
        }

        // Update the current layout if there are any layouts left
        if (newLayout) {
            set_layout(newLayout);
        } else {
            data.current_layout = null;
            data.chars = "";
        }

        save_layouts();
        render_layout();
    } else {
        alert("No layout selected or layout does not exist.");
    }
}


function reset_layouts() {
    layouts = { ...defaultLayouts };
    set_layout("qwerty");
    save_layouts();
    render_layout();
}


$(document).ready(function () {
    load_audio();
    load_layouts(); // Load saved layouts from localStorage
    if (localStorage.data != undefined) {
        load();
        render();
    }
    else {
        set_level(1);
    }
    $(document).keypress(keyHandler);
});

function start_stats() {
    start_time = start_time || Math.floor(new Date().getTime() / 1000);
}

function update_stats() {
    if (start_time) {
        var current_time = (Math.floor(new Date().getTime() / 1000));
        ratio = Math.floor(
            hits_correct / (hits_correct + hits_wrong) * 100
        );
        hpm = Math.floor(
            (hits_correct + hits_wrong) / (current_time - start_time) * 60
        );
        if (!isFinite(hpm)) { hpm = 0; }
    }
}


function set_level(l) {
    data.in_a_row = {};
    for (var i = 0; i < data.chars.length; i++) {
        data.in_a_row[data.chars[i]] = data.consecutive;
    }
    data.in_a_row[data.chars[l]] = 0;
    data.level = l;
    data.word_index = 0;
    data.word_errors = {};
    data.word = generate_word();
    data.keys_hit = "";
    save();
    render();
}

function set_layout(l) {
    data.current_layout = l
    data.chars = layouts[l]
    data.in_a_row = {};
    for (var i = 0; i < data.chars.length; i++) {
        data.in_a_row[data.chars[i]] = data.consecutive;
    }
    data.word_index = 0;
    data.word_errors = {};
    data.word = generate_word();
    data.keys_hit = "";
    save();
    render();
}


function keyHandler(e) {
    if (isModalOpen) {
        // Allow input to modal fields
        return;
    }
    start_stats();

    var key = String.fromCharCode(e.which);
    if (data.chars.indexOf(key) > -1) {
        e.preventDefault();
    }
    else {
        return;
    }
    data.keys_hit += key;
    if (key == data.word[data.word_index]) {
        hits_correct += 1;
        data.in_a_row[key] += 1;
        play_audio_sample("correct");
    }
    else {
        hits_wrong += 1;
        data.in_a_row[data.word[data.word_index]] = 0;
        data.in_a_row[key] = 0;
        play_audio_sample("mistake");
        data.word_errors[data.word_index] = true;
    }
    data.word_index += 1;
    if (data.word_index >= data.word.length) {
        setTimeout(next_word, 400);
    }

    update_stats();

    render();
    save();
}

function next_word() {
    if (get_training_chars().length == 0) {
        level_up();
    }
    data.word = generate_word();
    data.word_index = 0;
    data.keys_hit = "";
    data.word_errors = {};
    update_stats();

    render();
    save();
}


function level_up() {
    if (data.level + 1 <= data.chars.length - 1) {
        play_audio_sample("level_up");
    }
    l = Math.min(data.level + 1, data.chars.length);
    set_level(l);
}


function save() {
    localStorage.data = JSON.stringify(data);
}


function load() {
    data = JSON.parse(localStorage.data);
}


function load_audio() {
    audio.samples = {};
    audio.context = new (window.AudioContext || window.webkitAudioContext)();
    load_audio_sample("correct", "click.wav");
    load_audio_sample("mistake", "clack.wav");
    load_audio_sample("level_up", "ding.wav");
}


function load_audio_sample(name, url) {
    if (!audio.samples[name]) {
        // fetch the .wav file via XMLHttpRequest as jQuery doesn't support 'arraybuffer' dataType
        var request = new XMLHttpRequest();
        request.open("GET", url, true);
        request.responseType = "arraybuffer";
        request.onload = function () {
            audio.context.decodeAudioData(request.response).then(function (buffer) {
                audio.samples[name] = buffer;
            });
        };
        request.send();
    }
}


function play_audio_sample(name) {
    if (audio.samples[name]) {
        var source = audio.context.createBufferSource();
        source.buffer = audio.samples[name];
        source.onended
        source.connect(audio.context.destination);
        source.start();
    }
}


function add_layout() {
    var layout_name = document.getElementById('new_layout_name').value;
    var manual_entry_input = stringDivide(document.getElementById('manual_entry_input').value);
    var number_row_input = stringDivide(document.getElementById('number_row_input').value);
    var home_row_input = stringDivide(document.getElementById('home_row_input').value);
    var top_row_input = stringDivide(document.getElementById('top_row_input').value);
    var bottom_row_input = stringDivide(document.getElementById('bottom_row_input').value);
    var thumbs_input = stringDivide(document.getElementById('thumbs_input').value);
    var special_chars_input = stringDivide(document.getElementById('special_chars_input').value);
    var shifted_special_chars_input = stringDivide(document.getElementById('shifted_special_chars_input').value);
    var uppercase_letters_checkbox = document.getElementById('uppercase_letters_checkbox').checked;


    let layout_string;
    if (layout_name && manual_entry_input) {
        layout_string = " " + manual_entry_input;
    } else if (layout_name && (home_row_input || top_row_input || bottom_row_input || number_row_input || special_chars_input || uppercase_letters_input || shifted_special_chars_input || thumbs_input)) {
        // Combine the inputs to form the layout string
        layout_string = " " + home_row_input + top_row_input + bottom_row_input + thumbs_input;

        if (uppercase_letters_checkbox) {
            layout_string += home_row_input.toUpperCase() + top_row_input.toUpperCase() + bottom_row_input.toUpperCase() + thumbs_input.toUpperCase();
        }

        layout_string += number_row_input + special_chars_input + shifted_special_chars_input;
    } else {
        alert("Please add a name, at least one row, or a manual entry.");
        return;
    }

    layouts[layout_name] = layout_string;
    save_layouts(); // Save layouts to localStorage
    render_layout();
    closeModal();
}


function openModal() {
    document.getElementById('addLayoutModal').style.display = "block";
    isModalOpen = true;
}

function closeModal() {
    document.getElementById('addLayoutModal').style.display = "none";
    isModalOpen = false;
}

// Close the modal when the user clicks anywhere outside of the modal
window.onclick = function (event) {
    var modal = document.getElementById('addLayoutModal');
    if (event.target == modal) {
        modal.style.display = "none";
        isModalOpen = false;
    }
}

function render() {
    render_layout();
    render_level();
    render_word();
    render_level_bar();
    render_rigor();
    render_stats();
}
function render_layout() {
    var layouts_html = "<span id='layout'>";
    for (var layout in layouts) {
        if (data.current_layout == layout) {
            layouts_html += "<span style='color: #F00' onclick='set_layout(\"" + layout + "\");'> "
        } else {
            layouts_html += "<span style='color: #AAA' onclick='set_layout(\"" + layout + "\");'> "

        }
        layouts_html += layout + "</span>";

    }
    // Add button to open modal for adding new layouts
    layouts_html += `
        <div>
            <button onclick="openModal()">Add Layout</button>
				<button onclick="delete_current_layout()">Delete Current Layout</button>
				<button onclick="reset_layouts()">Reset Layouts to Default</button>
        </div>
    `;
    layouts_html += "</span>";

    $("#layout").html('Choose layout : ' + layouts_html);
}


function render_level() {
    var chars = "<span id='level-chars-wrap'>";
    var level_chars = get_level_chars();
    var training_chars = get_training_chars();
    for (var c in data.chars) {
        if (training_chars.indexOf(data.chars[c]) != -1) {
            chars += "<span style='color: #F00' onclick='set_level(" + c + ");'>"
        }
        else if (level_chars.indexOf(data.chars[c]) != -1) {
            chars += "<span style='color: #000' onclick='set_level(" + c + ");'>"
        }
        else {
            chars += "<span style='color: #AAA' onclick='set_level(" + c + ");'>"
        }
        if (data.chars[c] == ' ') {
            chars += "&#9141;";
        }
        else {
            chars += data.chars[c];
        }
        chars += "</span>";
    }
    chars += "</span>";
    $("#level-chars").html('click to set level: ' + chars);
}

function render_rigor() {
    chars = "<span id='rigor-number' onclick='inc_rigor();'>";
    chars += '' + data.consecutive;
    chars += '<span>';
    $('#rigor').html('click to set required repititions: ' + chars);
}

function render_stats() {
    $("#stats").text([
        "hits per minute: ", hpm, " ",
        "correctness: ", ratio, "%"
    ].join(""));
}

function inc_rigor() {
    data.consecutive += 1;
    if (data.consecutive > 9) {
        data.consecutive = 2;
    }
    render_rigor();
}


function render_level_bar() {
    training_chars = get_training_chars();
    if (training_chars.length == 0) {
        m = data.consecutive;
    }
    else {
        m = 1e100;
        for (c in training_chars) {
            m = Math.min(data.in_a_row[training_chars[c]], m);
        }
    }
    m = Math.floor($('#level-chars-wrap').innerWidth() * Math.min(1.0, m / data.consecutive));
    $('#next-level').css({ 'width': '' + m + 'px' });

}

function render_word() {
    var word = "";
    for (var i = 0; i < data.word.length; i++) {
        sclass = "normalChar";
        if (i > data.word_index) {
            sclass = "normalChar";
        }
        else if (i == data.word_index) {
            sclass = "currentChar";
        }
        else if (data.word_errors[i]) {
            sclass = "errorChar";
        }
        else {
            sclass = "goodChar";
        }
        word += "<span class='" + sclass + "'>";
        if (data.word[i] == " ") {
            word += "&#9141;"
        }
        else if (data.word[i] == "&") {
            word += "&amp;"
        }
        else {
            word += data.word[i];
        }
        word += "</span>";
    }
    var keys_hit = "<span class='keys-hit'>";
    for (var d in data.keys_hit) {
        if (data.keys_hit[d] == ' ') {
            keys_hit += "&#9141;";
        }
        else if (data.keys_hit[d] == '&') {
            keys_hit += "&amp;";
        }
        else {
            keys_hit += data.keys_hit[d];
        }
    }
    for (var i = data.word_index; i < data.word_length; i++) {
        keys_hit += "&nbsp;";
    }
    keys_hit += "</span>";
    $("#word").html(word + "<br>" + keys_hit);
}


function generate_word() {
    word = '';
    for (var i = 0; i < data.word_length; i++) {
        c = choose(get_training_chars());
        if (c != undefined && c != word[word.length - 1]) {
            word += c;
        }
        else {
            word += choose(get_level_chars());
        }
    }
    return word;
}


function get_level_chars() {
    return data.chars.slice(0, data.level + 1).split('');
}

function get_training_chars() {
    var training_chars = [];
    var level_chars = get_level_chars();
    for (var x in level_chars) {
        if (data.in_a_row[level_chars[x]] < data.consecutive) {
            training_chars.push(level_chars[x]);
        }
    }
    return training_chars;
}

function choose(a) {
    return a[Math.floor(Math.random() * a.length)];
}
