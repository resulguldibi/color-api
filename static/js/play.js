	googkeToken = window.localStorage.getItem("googleToken")
	googleLoginRequest = {};
	googleLoginRequest.token = googkeToken

	validateColorRequest = {}
	
	var helpedColors = [];

    initUser();
    
    function initUser(){
       var token =  window.localStorage.getItem("colorToken");
       if (token == null || token == undefined || token == ""){
            loginWithGoogle();
       }       

       prepareUserInfo();
    }    

	

    function prepareUserInfo(){
        var userInfo = getUserInfo();
        $('#imgUserPicture').css("display","none");
        if (userInfo != null && userInfo != undefined){
           $('#imgUserPicture').attr("src",userInfo.picture)
           $('#imgUserPicture').attr("title", userInfo.email)   
           $('#imgUserPicture').css("display", "block");
        }
    }
    
    function getUserInfo(){
       return parseJwt(window.localStorage.getItem("colorToken"))
    }

	function loginWithGoogle() {
	    var loginWithGoogleResponse;
	    $.ajax({
	        type: 'POST',
	        url: '../google/oauth2/token',
	        async: false,
	        contentType: 'application/json',
	        success: function (result) {
	            window.localStorage.setItem("colorToken", result.token);
                loginWithGoogleResponse = result;                
	        },
	        processData: false,
	        data: JSON.stringify(googleLoginRequest)
	    });

	    return loginWithGoogleResponse
	}

	function validateColors(validateColorRequest) {
	    var validateColorResponse;
	    $.ajax({
	        type: 'POST',
	        url: '../validate',
	        async: false,
	        contentType: 'application/json',
	        headers: {
	            'Authorization': window.localStorage.getItem("colorToken"),
	            'RaundKey': window.localStorage.getItem("raundKey")
	        },
	        success: function (result) {	            
                validateColorResponse = result
                setRaundStartPoint(result.raundPoint)
                setTotalPoint(result.totalPoint)
			},
			error: function (jqXHR, textStatus, errorThrown) {
				validateColorResponse = JSON.parse(jqXHR.responseText)
			},
	        processData: false,
	        data: JSON.stringify(validateColorRequest)
	    });

	    return validateColorResponse
	}

	function getColors() {
	    var getColorsResponse;


	    var selectedLevel = window.localStorage.getItem("selectedLevel");
	    var level = 2
	    if (selectedLevel != undefined && selectedLevel != null && selectedLevel != "") {
	        level = parseInt(selectedLevel)
	    } else {
	        window.localStorage.setItem("selectedLevel", level);
	        $('#slcLevels option[value=' + level + ']').first().attr('selected', 'selected')
	    }


	    $.ajax({
	        type: 'GET',
	        async: false,
	        url: '../colors?level=' + level,
	        headers: {
	            'Authorization': window.localStorage.getItem("colorToken")
	        },
	        success: function (result) {
	            
	            getColorsResponse = result
                window.localStorage.setItem("raundKey", result.code)
                setRaundStartPoint(result.raundStartPoint)
				setTotalPoint(result.totalPoint)
				helpedColors = []

	        },
	        processData: false
	    });

	    return getColorsResponse
	}

	function getLevels() {
	    var getLevelResponse
	    $.ajax({
	        type: 'GET',
	        url: '../levels',
	        async: false,
	        headers: {
	            'Authorization': window.localStorage.getItem("colorToken")
	        },
	        success: function (result) {
	            
	            getLevelResponse = result
	        },
	        processData: false
	    });

	    return getLevelResponse
	}

	function getRankings() {
	    var getRankingResponse
	    $.ajax({
	        type: 'GET',
	        url: '../rankings',
	        async: false,
	        headers: {
	            'Authorization': window.localStorage.getItem("colorToken")
	        },
	        success: function (result) {
	            
	            getRankingResponse = result
	        },
	        processData: false
	    });

	    return getRankingResponse
    }
    
    function getHelp() {

	    var getHelpResponse;

	    $.ajax({
	        type: 'GET',
	        async: false,
	        url: '../help',
	        headers: {
	            'Authorization': window.localStorage.getItem("colorToken"),
	            'RaundKey': window.localStorage.getItem("raundKey")
	        },
	        success: function (result) {
	            
	            getHelpResponse = result


	        },
	        processData: false
	    });

	    return getHelpResponse
    }

    function getStepHelp(selectedColors) {

        var getStepHelpResponse;

        $.ajax({
            type: 'GET',
            async: false,
            url: '../stephelp?colors=' + JSON.stringify(selectedColors),
            headers: {
                'Authorization': window.localStorage.getItem("colorToken"),
                'RaundKey': window.localStorage.getItem("raundKey")
            },
            success: function (result) {
                getStepHelpResponse = result
                setRaundStartPoint(result.point)
            },
            error: function (jqXHR, textStatus, errorThrown) {
                
               getStepHelpResponse = JSON.parse(jqXHR.responseText)
            },
            processData: false
        });

        return getStepHelpResponse
	}

	 function getRaundHistory() {

	 	var getRaundHistoryResponse;

	 	$.ajax({
	 		type: 'GET',
	 		async: false,
	 		url: '../history/raund',
	 		headers: {
	 			'Authorization': window.localStorage.getItem("colorToken"),
	 			'RaundKey': window.localStorage.getItem("raundKey")
	 		},
	 		success: function (result) {
	 			getRaundHistoryResponse = result
	 			debugger;
	 		},
	 		error: function (jqXHR, textStatus, errorThrown) {

	 			getRaundHistoryResponse = JSON.parse(jqXHR.responseText)
	 		},
	 		processData: false
	 	});

	 	return getRaundHistoryResponse
	 }
	


	//amac bir rengin hangi n rengin kar�s�m�ndan elde edildi�ini bulmak, bunun i�in toplam 5 * n + 1 adet random renk olusturulur.	
	//bu renkler generate edilirken, ilk olarak n adet random renk elde edilir. sonras�nda bu n rengin kar�s�mlar�ndan bir renk elde edilir.
	//geriye kalan (5 * n + 1 - (n +1)) adet random renk elde edilir. 
	//bu renklerden 1 tanesi kar�s�m sonucu bulunmas� gereken renk 5 * n tanesi ise, bu rengi elde etmek i�in kullan�labilecek renk seti olmal�d�r.	

	var selectedLevel = window.localStorage.getItem("selectedLevel");
	var n = 2;

	if (selectedLevel != undefined && selectedLevel != null && selectedLevel != "") {
	    n = parseInt(selectedLevel)
	    if (n <= 0 || n >= 5) {
	        n = 2;
	    }
	}

	$('#randomColors').attr("class", "level" + n + "RandomContainer");

	var totalColorCount = 5 * n + 1

	clearSettings();
	window.stepNumber = 0;
	initColors();
	initLevels();


    function clearSettings(){
        clearResult();
        $('#pnlRaundStartPoint').html('');
    }

    function clearResult(){
        $('#imgResultHappy').css("display", "none");
        $('#imgResultAngry').css("display", "none");
        $('#imgResultHelp').css("display", "none");
        $('#pnlTextResult').css("display", "none");
        $('#pnlTextResult').html('')
	}
	
	function displayHelpResult() {
		$('#imgResultHelp').css("display", "block");
	}

	function displayMessageResult(messsage) {
		$('#pnlTextResult').css("display", "block");
		$('#pnlTextResult').html(messsage)
	}

	 function setRaundStartPoint(point) {
	 	$('#pnlRaundStartPoint').html(point);
	 }

	 function setTotalPoint(point) {
	 	$('#pnlTotalPoint').html(point)
	 }


	 

    function stepHelp(){
        var selectedColors = getSelectedColors();
        var getStepHelpResponse = getStepHelp(selectedColors);

         var listContainer = $('.level' + n + 'Container');

         if (getStepHelpResponse != null && getStepHelpResponse != undefined && getStepHelpResponse.color != undefined && getStepHelpResponse.color !=null){

             for (var i = 0; i < listContainer.length; i++) {
                 var containerItem = listContainer[i]
                 var choiseItem = containerItem.childNodes[1];

                 var color = {}
                 color.r = parseInt(choiseItem.getAttribute("r"))
                 color.g = parseInt(choiseItem.getAttribute("g"))
                 color.b = parseInt(choiseItem.getAttribute("b"))
			   
				 var istried = choiseItem.getAttribute("istried");
                 var isselected = isColorEquals(getStepHelpResponse.color, color);

                 if (isselected || istried === 'true') {
					 
					if (isselected){
						if (!isColorExist(helpedColors, color)){
							helpedColors.push(color)							
						}	
						choiseItem.checked = false
					 }
					 if (istried === 'true' && isColorExist(helpedColors, color)) {
						choiseItem.checked = false
					 }
				                	
                 	onDivContainerClick(containerItem, false);
                 } 
             }

             var colors = getSelectedColors();
			 prepareMixColor(colors);
			 clearResult();
             displayHelpResult();
         } else if (!getStepHelpResponse.issuccess){
			  clearResult();
			  	displayHelpResult();
                displayMessageResult(getStepHelpResponse.message)
         }      
    }
	

	function help() {

	    var listContainer = $('.level' + n + 'Container');

	    var getHelpResponse = getHelp();

	    if (getHelpResponse != null && getHelpResponse != undefined && getHelpResponse.selectedColors != null && getHelpResponse.selectedColors != undefined && getHelpResponse.selectedColors.length > 0) {
	       
	        for (var i = 0; i < listContainer.length; i++) {
	            var containerItem = listContainer[i]
	            var choiseItem = containerItem.childNodes[1];

	            var color = {}
	            color.r = parseInt(choiseItem.getAttribute("r"))
	            color.g = parseInt(choiseItem.getAttribute("g"))
	            color.b = parseInt(choiseItem.getAttribute("b"))

	            var istried = choiseItem.getAttribute("istried");
	            var isselected = isColorExist(getHelpResponse.selectedColors, color);

	            if (isselected) {
	                choiseItem.checked = false
	                onDivContainerClick(containerItem, false);	                
	            } else if (istried === 'true') {
	                onDivContainerClick(containerItem, false);	                
	            }
            }
            
            var colors = getSelectedColors()
            prepareMixColor(colors)
			displayHelpResult();
			setRaundStartPoint(getHelpResponse.point)
	    }
    }
    
    

	function isColorExist(colors, color) {

	    var isExist = false;
	    if (colors != null && colors != undefined && colors.length > 0) {
	        for (var i = 0; i < colors.length; i++) {
	            var item = colors[i];

	            if (item != null && item != undefined && color != null && color != undefined) {

	                if (color.r === item.r && color.g === item.g && color.b === item.b) {
	                    isExist = true;
	                    break;
	                }
	            }
	        }
	    }

	    return isExist
    }
    

    function isColorEquals(color1, color2) {

        var isEquals = false;
        if (color1 != null && color1 != undefined && color2 != null && color2 != undefined) {
           if (color1.r === color2.r && color1.g === color2.g && color1.b === color2.b) {
               isEquals = true;               
           }
        }

        return isEquals
    }


	function onLevelChange(item) {
	    if (item.value == -1) {
	        return false;
	    }
	    n = item.value;
	    window.localStorage.setItem("selectedLevel", n)
	    totalColorCount = 5 * n + 1;

	    
	    clearSettings();
	    window.stepNumber = 0;

	    $('#randomColors').attr("class", "level" + n + "RandomContainer")

	    refreshItems();
    }
    
   

    function getSelectedColors(){
        var list = $('.choise'); //document.getElementsByClassName('choise');
        var colors = [];
        for (var i = 0; i < list.length; i++) {
            if (list[i].checked) {
                var item = list[i]
                var color = {}
                color.r = parseInt(item.getAttribute("r"))
                color.g = parseInt(item.getAttribute("g"))
                color.b = parseInt(item.getAttribute("b"))
                colors.push(color)
            }
        }

        return colors
    }

    function prepareMixColor(colors) {
        
        validateColorRequest.selectedColors = colors

        var mixColor = generateMixColor(colors);


        $('#mainColor').children('.mixColor').remove('.mixColor')


        var list = $('.actualColor').css("margin", "0px 10px 0px 0px");


        $('#mainColor').css("width", "224px").css("clear", "");

        var divMixColor = $("<div>");
        $(divMixColor).css("background", getColorHexCode(mixColor));

        $(divMixColor).attr("class", "mixColor");
        $(divMixColor).attr("r", mixColor.r);
        $(divMixColor).attr("g", mixColor.g);
        $(divMixColor).attr("b", mixColor.b);
        $(divMixColor).css("clear", "");

        $("#mainColor").append(divMixColor)

        return mixColor
    }

	function Validate(mixColor, colors) {

	    
	    clearSettings();

	    var validateColorRequest = {}
	   
	    var actualColorList = $('.actualColor');
	    var actualColorItem = actualColorList[0]

	    actualColor = {}
	    actualColor.r = parseInt(actualColorItem.getAttribute("r"))
	    actualColor.g = parseInt(actualColorItem.getAttribute("g"))
	    actualColor.b = parseInt(actualColorItem.getAttribute("b"))

	    window.stepNumber += 1;   

	    if (mixColor != null && mixColor != undefined && actualColor != null && actualColor != undefined) {
	        validateColorRequest.selectedColors = colors
	        var validateColorsResponse = validateColors(validateColorRequest)
	        
	        if (validateColorsResponse != null && validateColorsResponse != undefined) {
                 clearResult()	            
				if (validateColorsResponse.isValid) {
					$('#imgResultHappy').css("display", "block");
				} else if (!validateColorsResponse.issuccess) {							
					$('#imgResultAngry').css("display", "block");
					displayMessageResult(validateColorsResponse.message)
				}
	        }
	    }
	}


	function prepareSelectedDivContainer(item) {
	    if (item.childNodes[1].checked) {
	        item.style.margin = '-30px 10px 0px 0px';
	        item.childNodes[1].setAttribute('istried', 'true');
	    } else {
	        item.style.margin = '0px 10px 0px 0px';
	        item.childNodes[1].setAttribute('istried', 'false');
	    }
	}

	function onDivContainerClick(item, validate) {
        clearResult();
	    item.childNodes[1].checked = !item.childNodes[1].checked;

	    prepareSelectedDivContainer(item);

	    var list = $('.choise');
	    var triedItemCount = 0;
	    for (var i = 0; i < list.length; i++) {
	        var choiseItem = list[i]

	        var istried = choiseItem.getAttribute("istried");

	        if (istried === 'true') {
	            triedItemCount = triedItemCount + 1;
	        }
        }

        var colors = getSelectedColors()
        var mixColor = prepareMixColor(colors)

	    if (triedItemCount == n) {

	        var listContainer = $('.level' + n + 'Container');

	        for (var i = 0; i < listContainer.length; i++) {
	            var containerItem = listContainer[i]
	            var choiseItem = containerItem.childNodes[1];
	            var istried = choiseItem.getAttribute("istried");

	            if (istried === 'false') {
	                containerItem.style.pointerEvents = 'none';
	                containerItem.style.opacity = '0.4';
	            }
	        }
           
	        if (validate) {
	            Validate(mixColor, colors);
	        }
	    } else {

	        var listContainer = $('.level' + n + 'Container');

	        for (var i = 0; i < listContainer.length; i++) {
	            var containerItem = listContainer[i]
	            containerItem.style.pointerEvents = 'all';
	            containerItem.style.opacity = '1';
	        }
	    }

	}


	function initLevels() {
	    var getLevelsResponse = getLevels();

	    $('#slcLevels').append($('<option>', {
	        value: "-1",
	        text: 'level'
	    }));

	    if (getLevelsResponse != null && getLevelsResponse != undefined && getLevelsResponse.levelCount > 0) {

	        for (var i = 1; i <= getLevelsResponse.levelCount; i++) {
	            $('#slcLevels').append($('<option>', {
	                value: i,
	                text: i
	            }));
	        }

	        var selectedLevel = window.localStorage.getItem("selectedLevel");
	        var level = getLevelsResponse.defaultLevel

	        if (selectedLevel != undefined && selectedLevel != null && selectedLevel != "") {
	            level = parseInt(selectedLevel)
	        } else {
	            window.localStorage.setItem("selectedLevel", getLevelsResponse.defaultLevel)
	        }

	        $('#slcLevels option[value=' + level + ']').first().attr('selected', 'selected')
	    }
	}

	function initColors() {

	    
	    clearSettings();

	    var getColorsResponse = getColors();
	    var mixColor = getColorsResponse.mixedColor

	    var finalColors = getColorsResponse.randomColors;

	    var divMainColor = $('<div>');
	    $(divMainColor).css("background", getColorHexCode(mixColor))


	    $(divMainColor).attr("class", "actualColor");

	    $(divMainColor).attr("r", mixColor.r);
	    $(divMainColor).attr("g", mixColor.g);
        $(divMainColor).attr("b", mixColor.b);
        $(divMainColor).attr("title", mixColor.name);

        $('#mainColor').append(divMainColor)
        //$('#mainColor').append($('<div>').attr("class", "pnlColorName").html(mixColor.name))

	    for (var i = 0; i < finalColors.length; i++) {
	        var div = $('<div>');
	        var divContainer = $('<div>');

	        $(divContainer).attr("onclick", "onDivContainerClick(this,true)");
	        $(divContainer).attr("class", "level" + n + "Container");

            $(div).attr("class", "level" + n + "RandomColor");
            $(div).attr("title", finalColors[i].name);
	        $(div).css("background", getColorHexCode(finalColors[i]))
            $(divContainer).append(div)
            //$(divContainer).append($('<div>').attr("class", "pnlColorName").html(finalColors[i].name))

	        var choise = $("<input>");

	        $(choise).css("cursor", "pointer");

	        $(choise).attr("type", "checkbox");

	        $(choise).attr("r", finalColors[i].r);
	        $(choise).attr("g", finalColors[i].g);
	        $(choise).attr("b", finalColors[i].b);

	        $(choise).attr("isselected", finalColors[i].isselected);

	        $(choise).attr("class", "choise");
	        $(choise).attr("istried", "false");

	        $(choise).css("display", "none");
	        $(choise).css("width", "100px");
	        $(choise).css("margin", "5px 0px 0px 0px");
	        $(divContainer).append(choise)

	        var divSeperator = $("<div>");
	        $(divSeperator).css("clear", "both");
	        $(divContainer).append(divSeperator)

	        $('#randomColors').append(divContainer)
	    }
	}

	function refreshItems() {

	    window.stepNumber = 0
	    clearChilds("mainColor")
	    clearChilds("randomColors")

	    var mainColor = $('#mainColor');
	    $(mainColor).css("width", "112px");


	    initColors()
	}

	function clearChilds(element) {

	    var myNode = document.getElementById(element);
	    while (myNode.firstChild) {
	        myNode.removeChild(myNode.firstChild);
	    }
	}



	function getColorHexCode(color) {
	    return "#" + getColor(color.r) + getColor(color.g) + getColor(color.b);
	}

	function generateMixColor(colors) {
	    var r = 0,
	        g = 0,
	        b = 0;
	    for (var i = 0; i < colors.length; i++) {
	        r = r + colors[i].r
	        g = g + colors[i].g
	        b = b + colors[i].b
	    }

	    var color = {}
	    color.r = Math.floor(r / colors.length)
	    color.g = Math.floor(g / colors.length)
	    color.b = Math.floor(b / colors.length)

	    return color
	}

	function generateRandomColor() {
	    var color = {}
	    color.r = generateRandomNumer(256)
	    color.g = generateRandomNumer(256)
	    color.b = generateRandomNumer(256)
	    return color
	}

	function getColor(code) {
	    if (code < 16)
	        return code.toString(16) + code.toString(16)
	    else
	        return code.toString(16)
	}

	function generateRandomNumer(maxNumber) {
	    return Math.floor(Math.random() * maxNumber);
    }
    
    function parseJwt(token) {
        var base64Url = token.split('.')[1];        
        var base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        return JSON.parse(window.atob(base64));
    }

    