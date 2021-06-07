$(document).ready(function(){
    loadList(0,'');
    loadGenres();
    loadModes();     
    $("body").delegate("a.genresclick","click",function () {
        var uid = $(this).attr('data-id');
        loadList(1,uid);
        jQuery('#backbutton').show();
       });

       $("body").delegate("a#favoritesclick","click",function () {
        loadList(3,'');
        jQuery('#backbutton').show();
       });

       

       $("body").delegate("a#backpress","click",function () {
        loadList(0,'');
        jQuery('#backbutton').hide();
       });

       $("body").delegate("button#sideload","click",function () {
        uid = $(this).attr('data-id');
        $(this).text("Sideloading...");
        $(this).prop('disabled', true);
        $(this).attr("class","btn btn-warning my-lg-0");
        jQuery.ajax({
            type: "GET",
            url: "http://127.0.0.1:12346/sideload?item="+uid,
            dataType: "json",
            success: function (data) {
                if (data.Type > 0) {
                    console.log("ok");
                }             
            },
            error: function () {
                console.log("request no. 7 have failed");
                jQuery('#sideload').text("Sideload to legacy pc");
                jQuery('#sideload').attr("class","btn btn-danger my-lg-0");
                jQuery("#sideload").prop("disabled",false);
            }
        }).then(function () {
            jQuery('#sideload').text("Sideload to legacy pc");
            jQuery('#sideload').attr("class","btn btn-danger my-lg-0");
            jQuery("#sideload").prop("disabled",false);
        });
       });

       $("body").delegate("button#openbtn","click",function () {
        uid = $(this).attr('data-id');
        $(this).text("Running...");
        $(this).prop('disabled', true);
        $(this).attr("class","btn btn-warning my-lg-0");
        jQuery.ajax({
            type: "GET",
            url: "http://127.0.0.1:12346/open?item="+uid,
            dataType: "json",
            success: function (data) {
                if (data.Type > 0) {
                    jQuery('#openbtn').attr("class","btn btn-success my-lg-0");
                    jQuery("#openbtn").text("Play");
                } else {
                    jQuery('#openbtn').attr("class","btn btn-info my-lg-0");
                    jQuery("#openbtn").text("Download & Play");
                }
                             
            },
            error: function () {
                console.log("request no. 6 have failed ");
            }
        }).then(function () {
            jQuery("#openbtn").prop("disabled",false);
        });
       });

       $("body").delegate("button#plusbtn","click",function () {
        uid = $(this).attr('data-id');
        jQuery.ajax({
            type: "GET",
            url: "http://127.0.0.1:12346/fav?item="+uid,
            dataType: "json",
            success: function (data) {
                if (data.Fav) {
                    jQuery('#plusbtn').attr("class","btn btn-info my-lg-0");
                    jQuery("#plusbtn").text("-");
                } else {
                    // need to remove from list
                    jQuery('#favitem_'+uid).remove();
                    jQuery('#plusbtn').attr("class","btn btn-success my-lg-0");
                    jQuery("#plusbtn").text("+");
                }             
            },
            error: function () {
                console.log("request no. 4 have failed ");
            }
        }).then(function () {
            //jQuery('#puslapis_load').hide();
        });
       });

       $("body").delegate("a.modesclick","click",function () {
        var uid = $(this).attr('data-id');
        loadList(2,uid);
        jQuery('#backbutton').show();
       });

       $("body").delegate("button#searchbutton","click",function () {
        var searchas = $('#searchfield').val();
        loadList(0,searchas)
        jQuery('#backbutton').show();
       });

       $("#searchfield").on( "keydown", function(event) {
        if(event.which == 13) {
        var searchas = $('#searchfield').val();
        loadList(0,searchas)
        jQuery('#backbutton').show();
        }
      });
    
        $("body").delegate("a.gamelistclick","click",function () {
            var uid = $(this).attr('data-id');
            jQuery.ajax({
                type: "GET",
                url: "http://127.0.0.1:12346/gamedetails?uid="+uid,
                dataType: "json",
                success: function (data) {
                  jQuery('#gameheader').html(data.Name);
                  jQuery('#brief').html(data.Brief);
                  jQuery('#openbtn').attr('data-id',uid);
                  jQuery('#plusbtn').attr('data-id',uid);
                  jQuery('#sideload').attr('data-id',uid)
                  if (data.Installed) {
                      jQuery('#openbtn').attr("class","btn btn-success my-lg-0");
                      jQuery("#openbtn").text("Play");
                  } else {
                    jQuery('#openbtn').attr("class","btn btn-info my-lg-0");
                    jQuery("#openbtn").text("Download & Play");
                  }
                  if (data.Fav) {
                    jQuery('#plusbtn').attr("class","btn btn-info my-lg-0");
                    jQuery("#plusbtn").text("-");
                  } else {
                    jQuery('#plusbtn').attr("class","btn btn-success my-lg-0");
                    jQuery("#plusbtn").text("+");
                  }
                  if (data.Experimental) {
                    jQuery('#sideload').attr("class","btn btn-danger my-lg-0");
                    jQuery("#sideload").show();
                  }
                  jQuery("#plusbtn").show();
                  jQuery("#openbtn").show();

                  jQuery('#images').html("");
                  data.ImgUrls.forEach(function(item) {
                      jQuery('#images').append('<img src="'+abandonware_url+'/image.php?file='+item+'"/>');
                  });
                  jQuery('#videos').html("");
                  data.VideoUrls.forEach(function(item) {
                    var myId = getYoutubeId(item); 
                    jQuery('#videos').append('<iframe width="560" height="315" src="//www.youtube.com/embed/' + myId + '" frameborder="0" allowfullscreen></iframe>');
                });
                  jQuery("#descrlist").html('<li class="list-group-item d-flex justify-content-between align-items-center">Producer: '+data.Producer+'</li>');
                  jQuery("#descrlist").append('<li class="list-group-item d-flex justify-content-between align-items-center">Publisher: '+data.Publisher+'</li>');
                  jQuery("#descrlist").append('<li class="list-group-item d-flex justify-content-between align-items-center">Year: '+data.Year+'</li>');
                  var Genres = "";
                  var genres_count = 0;
                  data.Genres.forEach(function(item) {
                    if (genres_count == 0) {
                        Genres = item;
                    }
                    if (genres_count > 0) {
                        Genres = Genres + " "+ item;
                    }
                    genres_count++;
                });
                jQuery("#descrlist").append('<li class="list-group-item d-flex justify-content-between align-items-center">Genres: '+Genres+'</li>');
                var Modes = "";
                var modes_count = 0;
                data.Modes.forEach(function(item) {
                  if (modes_count == 0) {
                      Modes = item;
                  }
                  if (modes_count > 0) {
                      Modes = Modes + " "+ item;
                  }
                  modes_count++;
              });
                jQuery("#descrlist").append('<li class="list-group-item d-flex justify-content-between align-items-center">Modes: '+Modes+'</li>');
                jQuery('#descrlist').show();
                jQuery('#description').html(data.Description);
                },
                error: function () {
                    console.log("request no. 3 have failed " + url);
                }
            }).then(function () {
                //jQuery('#puslapis_load').hide();
            });
           });
   
            

})

function getYoutubeId(url) {
    var regExp = /^.*(youtu.be\/|v\/|u\/\w\/|embed\/|watch\?v=|\&v=)([^#\&\?]*).*/;
    var match = url.match(regExp);

    if (match && match[2].length == 11) {
        return match[2];
    } else {
        return 'error';
    }
}

var abandonware_url = "https://abandonware.club/";
var local_url = "http://127.0.0.1:12346/";

function loadGenres() {
    jQuery.ajax({
        type: "GET",
        url: local_url+"genreslist",
        dataType: "json",
        success: function (data) {
          jQuery('#genreslist').html("<br>");
          var genrecount = 0;
          data.forEach(function(item) {
            genrecount++;  
            jQuery('#genreslist').append('<a href="#" class="card-link genresclick" data-id="'+item.id+'">'+item.name+'</a> ');
            if (genrecount > 5) {
                jQuery("#genrelist").append('<br>');
            }
        });
           jQuery('#genres').append("<br><br>");
        },
        error: function () {
            console.log("request no. 1 have failed " + url);
        }
    }).then(function () {
        //jQuery('#puslapis_load').hide();
    });
}

function loadModes() {
    jQuery.ajax({
        type: "GET",
        url: local_url+"modeslist",
        dataType: "json",
        success: function (data) {
          jQuery('#modeslist').html("<br>");
          var genrecount = 0;
          data.forEach(function(item) {
            genrecount++;  
            jQuery('#modeslist').append('<a href="#" class="card-link modesclick" data-id="'+item.id+'">'+item.name+'</a> ');
            if (genrecount > 5) {
                jQuery("#modeslist").append('<br>');
            }
        });
           jQuery('#modeslist').append("<br><br>");
        },
        error: function () {
            console.log("request no. 2 have failed " + url);
        }
    }).then(function () {
        //jQuery('#puslapis_load').hide();
    });
}

function loadList(mode,searchstr) {
    jQuery.ajax({
        type: "GET",
        url: "http://127.0.0.1:12346/gamelist?mode="+mode+"&key="+searchstr,
        dataType: "json",
        success: function (data) {
          jQuery('#gamelistas').html("");
          data.forEach(function(item) {
              var idsas = "";
              if (mode == 3) {
              idsas = 'id="favitem_'+item.uid+'" ';
              }
            jQuery('#gamelistas').append('<a href="#" '+idsas+'class="list-group-item list-group-item-action gamelistclick" data-id="'+item.uid+'">'+item.name+'</a>');
        });
          
        },
        error: function () {
            console.log("request no. 4 have failed " + url);
        }
    }).then(function () {
        //jQuery('#puslapis_load').hide();
    });
}



