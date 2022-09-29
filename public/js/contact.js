// let nama = "rahmat";
// let umur = 26

// console.log(`nama saya ${nama} umur saya ${umur} tahun`);
// document.write(`nama saya ${nama} umur saya ${umur} tahun`)
// alert(`nama saya ${nama} umur saya ${umur} tahun`)

// operator
// let a = "aku"
// let b = "kamu"
// let result = a + b
// console.log(result);

// condition 
// if (result == "akukamu"){result="bahagia"} else {result = "tidak bahagia"}
// console.log(result);

// let nilai = 75

// if(nilai >= 75 ){
//     console.log("lulus");
// } else {
//     console.log("tidak lulus");
// }




function submitData (){

let name = document.getElementById("name").value
let email = document.getElementById("email").value
let phoneNumber = document.getElementById("phone-number").value
let subject = document.getElementById ("subject").value
let message = document.getElementById ("message").value


if ( name == ""||
email == ""||
phoneNumber == ""||
subject == ""||
message == "")
{
    return alert("Please input your information completely")} 
    else {
    alert("thank you, and hit OK to continue")
    }
    
// console.log(name);
// console.log(email);
// console.log(phoneNumber);
// console.log(subject);
// console.log(message);

// This is another mailto link:
/* <a href="mailto:someone@example.com?&subject=Summer Party&body=You%100are invited to a big summer party!" target="_top">Send mail!</a> */

let emailReciever = "Rahmat9654@gmail.com"  
let a = document.createElement(`A`)
// a.href=`mailto:${emailReciever}?&subject=${subject}&body=Hello my name is ${name}, ${message}, hit me up! my phone number is ${phoneNumber} and my email is ${email}.` 
a.setAttribute("href",`mailto:${emailReciever}?&subject=${subject}&body=Hello my name is ${name}, ${message}, hit me up! my phone number is ${phoneNumber} and my email is ${email}.`)
a.click()

}