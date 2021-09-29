
//// A simple function.
//function hello(longName) {
//  alert('Hello, ' + longName);
//}
//
//hello('New User');
//

//function get_date() {
//	ll = document.getElementsByClassName("date");
//	q = ll.item(0).innerHTML;
//	//console.log(q);
//	//qr = q.replace("<br>", "\n");
//	//z = q.split("<br>", -1);
//	z = q.split('<br>', -1);
//	for(i=0; i<z.length; z++) {
//		console.log(String(z[i]));
//	}
//}
//
//get_date();
//
function find_weekday(val) {
	f = val.indexOf("（");
	if(f > -1) {
		b = val.indexOf("）");
		if(b > -1) {
			return val.slice(f+1, b);
		}
	}
	return "";
}

function find_time(val) {
	f = val.indexOf("）");
	if(f > -1) {
		ret = val.slice(f+1).replace("：", ":");
		return ret;
	}
	return "";
}

function get_title() {
	ll = document.getElementsByClassName("ttl01");
    title = ll.item(0).innerHTML;
	//console.log(q);
	return title;
}

function find_time_line(v) {
	for(i=1; i<v.length; i++){
		pos = z[i].indexOf("：");
		if(pos > 0){
			return i;
		}
	}
	return -1;
}

function get_date() {
    ll = document.getElementsByClassName("date");
    q = ll.item(0).innerHTML;
    z = q.split('<br>', -1);
	pos = find_time_line(z);
	week_day = find_weekday(z[pos]);
	//console.log(week_day);
	time = find_time(z[pos]);
	//console.log(time);
	return [week_day, time];
}

function gen_json_queue() {
	station = "AT-X";
	cs = "true";
	const ret_json = `		"%%%TITLE%%%": {
			"start_time": "%%%TIME%%%",
			"weekday": "%%%WEEKDAY%%%",
			"station": "%%%STATION%%%",
			"is_cs": %%%CS%%%
		},`;
	[week_day, time] = get_date();
	title = get_title();
	r = ret_json.replace("%%%TITLE%%%", title);
	r = r.replace("%%%TIME%%%", time);
	r = r.replace("%%%WEEKDAY%%%", week_day);
	r = r.replace("%%%STATION%%%", station);
	r = r.replace("%%%CS%%%", cs);
	console.log(r);
	return r;
}

function copyToClipboard(text) {
  return navigator.clipboard.writeText(text).then(function() {
    alert('コピーしました')
  }).catch(function(error) {
    alert((error && error.message) || 'コピーに失敗しました')
  })
}

//get_date();
//get_title();

jjj = gen_json_queue();
copyToClipboard(jjj);
