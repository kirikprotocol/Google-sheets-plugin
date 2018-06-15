# Google Sheets Plugin
## Goal

To enable interworking between a MiniApps bot and Google Sheets whereby the bot uses a Google Sheets document as a database. It should make possible e.g. registration of users and recording their data into the table. The plugin also can:

- retrieve the last record from the table;
- notify administrators by their phone numbers (in Telegram) and email addresses;
- localize answers the system gives to users coming from various places.

## Usage

To enable interworking between the plugin and your Google Sheet document allow miniapps@miniappstesterbot.iam.gserviceaccount.com to edit it. To do it:

a. Create the sheet and press the SHARE button in the upper right corner:
![](https://i.imgur.com/lEjAFXd.png)
b. Then insert the above shown address in the People entry box and press Send button:
![](https://i.imgur.com/X8rP7vb.png)

Then you should open MiniApps Visual Editor and configure your Google Sheets bot (refer to [Constructing Google Sheets Bot](#constructing-google-sheets-bot)) and frame an appropriate request that the bot should send to the Google Sheets plugin (refer to [Framing Request to Google Sheets Plugin](#framing-request-to-google-sheets-instruction)).

## Constructing Google Sheets Bot

This chapter explains how to configure your bot using the Visual Editor to make it interoperable with Google Sheets Plugin.

### General Information About Visual Editor Configuration

* Connection of page names and spreadsheet columns.
When working with pages in the Visual Editor you can name them by setting their IDs:
![](https://i.imgur.com/CoCewNh.png)
![](https://i.imgur.com/wlm427f.png) ![](https://i.imgur.com/HICDpii.png)

The page ID will be the same as the name of its associated spreadsheet column:
![](https://i.imgur.com/cBRLT01.png)

The information contained in this column cells will be the recorded answers that users give on the associated page.
* Button IDs in the Visual Editor and their corresponding values in the column, answer evaluation (parameter evaluable 1 - evaluate).
Button IDs set in the Visual Editor are recorded to cells of the spreadsheet when these buttons are pushed.
![](https://i.imgur.com/93zkqLu.png) 
![](https://i.imgur.com/EocsJss.png)

* If keyboard input is needed (e.g. the user's email), you can set the default target page:
![](https://i.imgur.com/8jvu6Qj.png)
![](https://i.imgur.com/HTTYSnL.png)

The answer is evaluated by the button number and the corresponding score:
![](https://i.imgur.com/VlALsut.png)

If the answer is right, the user scores 1 point, if the answer is wrong - 0 points.

* Calling Plugin:

1. Create a page and change its type to External Service.
![](https://i.imgur.com/Xs2ByBm.png)
![](https://i.imgur.com/Yed4RxL.png)

2. Insert the URL that you have configured following the [Framing Request to Google Sheets instruction](#framing-request-to-google-sheets-instruction).
![](https://i.imgur.com/nwNerHB.png)

3. Switch on Transfer user answers, if you want to use the function.

4. If needed, add the callback parameter and name it "callback":
![](https://i.imgur.com/xEt0sEA.png)
![](https://i.imgur.com/VWoEIRh.png)

5. Choose the page which callback should lead to:
![](https://i.imgur.com/owTVsKP.png)
![](https://i.imgur.com/5JNxJeT.png)
        
 ## Framing Request to Google Sheets instruction

First thing to do is to form the main part of the request. To do it you need the ID of your Google Sheet:

1. Copy the sheet address from the browser:
![](https://i.imgur.com/HytWfPz.png)

2. Select the part in the middle of the address flanked by "/" on both sides, like this: https://docs.google.com/spreadsheets/d/**1GCXT5ii2NJxok6hpnjAXp3RQd6H_9TQs4pkKB6PDbZc**/edit?pli=1#gid=2147. This is the ID. Copy it.

The plugin has two main functions:
* Add a row to the sheet (its address is simply /)
* Retrieve the last row belonging to a particular user (its address is /**getLast**/)

Choose the function that you need and place its address instead of <function_address> in 
```
http://plugins.miniapps.run/MiniappsTesterGoogleSheetsBot<function_address> 
```
Add the sheet ID to: 
```
http://plugins.miniapps.run/MiniappsTesterGoogleSheetsBot<function_address>?spreadsheetId=
```

Which will result in an address line of this kind: 
```
http://plugins.miniapps.run/MiniappsTesterGoogleSheetsBot/?spreadsheetId=1GCXT5ii2NJxok6hpnjAXp3RQd6H_9TQs4pkKB6PDbZc
```

The next steps depend on your desired result. You can add a parameter:
```
&<parameter_name>=<parameter_value>,
```

where <parameter_name> – is the name of the parameter, and <parameter_value> is its value.

For instance, http://plugins.miniapps.run/MiniappsTesterGoogleSheetsBot/?spreadsheetId=1GCXT5ii2NJxok6hpnjAXp3RQd6H_9TQs4pkKB6PDbZc&evaluable=0&dispatch=1&sendEmail=1

Below is the list of the parameters and description of their values, added with possible tasks, to which they are applicable. *For more information on the tasks refer to [Constructing Google Sheets Bot](#constructing-google-sheets-bot).*

Parameter       |Function with which it is passed       |Mandatory      |Tasks           |Value          |
----------------|---------------------------------------|---------------|----------------|---------------|
spreadsheetId   |/ and /getLast/        |Yes            |see [Constructing Google Sheets Bot](#constructing-google-sheets-bot)  |Google Sheet ID        |
evaluable	|/                      |No (0 by default)|see [Constructing Google Sheets Bot](#constructing-google-sheets-bot)        |Whether or not user's answers should be evaluated(yes, if the user gains a score)      |
callback	|/ and /getLast/        |No (plugin returns its answer by default)|see [Constructing Google Sheets Bot](#constructing-google-sheets-bot)|Callback url (the address the dialog should be forwarded to after the plugin finishes its work)[1] |
translationTableTitle                   |/ and /getLast/|No (Translation by default)|see [Constructing Google Sheets Bot](#constructing-google-sheets-bot) |The name of the tab in the table that contains the translation[2]|
parameters	|/getLast/	        |Yes	        |see [Constructing Google Sheets Bot](#constructing-google-sheets-bot) |Comma separated parameter names that should be given to the user|
dispatch        |/              	|No (0 by default)|see [Constructing Google Sheets Bot](#constructing-google-sheets-bot)        |This parameter determines whether administrators specified in userTableTitle should be notified|
sendEmail       |/                      |No (0 by default)|see [Constructing Google Sheets Bot](#constructing-google-sheets-bot)        |This parameter determines whether administrators should be notified by email (used in conjunction with *dispatch*)
userTableTitle	|/                      |No (DispatchPhoneList by default)|see [Constructing Google Sheets Bot](#constructing-google-sheets-bot)|The name of the tab in the table that contains administrators' data**|

[1] - callback is entered in the editor like this:
![](https://i.imgur.com/Ngi2Y0G.png)

[2] - you should create a tab, name and format it accordingly in the table that is passed in spreadsheetId:

- Name: translationTableTitle; This tab contains translation of field names and bot response prefix. Format:
![](https://i.imgur.com/B7I1iyu.png)
(prefix – the standard parameter setting the beginning of the bot response. The rest is the parameters from the table. en and ru - language codes).

- Name: userTableTitle; This tab contains the administrators' data. Format:
![](https://i.imgur.com/E5OWaI3.png)
