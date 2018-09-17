// Code generated by go generate; DO NOT EDIT.
// +build !js

package Front


var HTMLAccount = HTMLPAGE(`<main application="account">

    <section page="tos">
        <div class="middle">

            <label>LICENSE</label>
            <tos class="tos">MIT License
Copyright (c) 2018 Inkeliz

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
associated documentation files (the "Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the
following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial
portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT
LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN
NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
            </tos>
            <input class="accept" type="submit" value="ACCEPT AND CONTINUE">

        </div>
    </section>

    <section page="index">
        <div class="middle">

            <input class="genSeed" type="submit" value="CREATE NEW WALLET">

            <hr>

            <input class="importSeed" type="submit" value="IMPORT SEEDFY">

        </div>
    </section>

    <section page="generate">
        <div class="middle">
            <alert class="invisible"></alert>

            <label>COPY YOUR SEEDFY</label>
            <textarea class="seed" rows="3" spellcheck="false"></textarea>

            <input class="continue" type="submit" value="CONTINUE">

            <alert>
                <icon class="icon-alert"/>
                <text>Losing your SEEDFY, or even the password, prevents access the wallet.</text>
            </alert>

        </div>
    </section>

    <section page="import">
        <div class="middle">
            <alert class="invisible"></alert>

            <label>ENTER YOUR SEEDFY</label>
            <textarea class="seed" rows="3" spellcheck="false"></textarea>

            <input class="continue" type="submit" value="CONTINUE">

            <alert>
                <icon class="icon-alert"/>
                <text>Keep the SEEDFY in a safe place, don't throw it away after access your wallet.</text>
            </alert>

        </div>
    </section>

    <section page="password">
        <div class="middle">
            <alert class="invisible"></alert>

            <label>PASSWORD</label>
            <input class="password" type="password">

            <label>2FA</label>
            <label class="checkbox"><input class="ask2fa" type="checkbox" value="true">Ask for two factor authentication</label>

            <input class="continue" type="submit" value="CONTINUE">

            <alert>
                <icon class="icon-alert"/>
                <text>You must remember your password in order to access the same wallet.</text>
            </alert>
        </div>
    </section>

    <section page="mfa">
        <div class="middle">

            <label>SCAN WITH YOUR NANOLLET 2FA</label>
            <div class="qrcode"></div>

            <input class="continue" type="submit" disabled value="CONTINUE">

        </div>
    </section>

    <section page="address">
        <div class="middle">

            <label>CHOOSE YOUR ADDRESS</label>
            <div>
                <button class="next">
                    <icon class="icon-next"/>
                </button>
                <button class="previous">
                    <icon class="icon-previous"/>
                </button>
                <select type="list" class="address box">
                </select>
            </div>

            <input class="continue" type="submit" value="CONTINUE">

        </div>
    </section>


</main>
`)
var HTMLBase = HTMLPAGE(`<html window-frame="solid" window-blurbehind="light" resizeable>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
    <title>Nanollet</title>

    <link rel="stylesheet" type="text/css" href="../css/style.css">
</head>
<body class="notconnected">


<header role="window-caption">
    <section class="logo">
        NANOLLET
    </section>
    <section class="account">
        <div class="box">
            <icon class="icon-coins"/>
            <div class="ammount"> </div>
        </div>
        <div class="box">
            <icon class="icon-nodes"/>
            <div class="nodes"> </div>
        </div>
        <button role="window-close">
            <icon id="end" class="icon-exit"></icon>
        </button>
    </section>
</header>

<section class="control">
</section>

<section class="dynamic">
</section>

<section class="notification">
    <!--<button class="notification"></button>-->
</section>
<!--
    Custom Input - DropZone:
    This input is used in "Nanofy" to select the file, it
 -->
<script type="text/tiscript">
    class DropZone : Element {
      event click() {
         var fn = view.selectFile(#open, "All Files (*.*)|*.*");
         if( fn ) {
          this.parent.select(".filepath").value = fn
          this.select(".name").value = fn.split("/")[fn.split("/").length-1]
          return true;
         }
      }
      event drag (evt) {
        if(evt.draggingDataType == #file) {
          return true;
        }
      }
      event drop (evt) {
        if(evt.draggingDataType == #file) {
          this.parent.select(".filepath").value = evt.dragging
          this.select(".name").value = evt.dragging.replace("\\", "/").split("/")[evt.dragging.replace("\\", "/").split("/").length-1]
          return true;
        }
      }
    }
</script>
</body>
</html>`)
var HTMLNanofy = HTMLPAGE(`<main application="nanofy">

    <section page="sign">
        <div class="middle">

            <label for="file">FILE</label>
            <dropzone type="file" name="file"><div class="content name">Drop the file here</div> </dropzone>
            <textarea type="hidden" class="filepath"></textarea>

            <input class="continue" type="submit" value="SIGN">

        </div>
    </section>

    <section page="verify">
        <div class="middle">

            <label for="file">FILE</label>
            <dropzone type="file" name="file"><div class="content name">Drop the file here</div> </dropzone>
            <textarea type="hidden" class="filepath"></textarea>

            <label for="address">ADDRESS</label>
            <textarea spellcheck="false" class="address" rows="2"></textarea>

            <input class="continue" type="submit" value="VERIFY">

        </div>
    </section>

</main>
`)
var HTMLNanollet = HTMLPAGE(`<main application="nanollet">

    <section page="send">
        <div class="middle">

            <label for="address">ADDRESS OR OPENCAP ALIAS</label>
            <textarea spellcheck="false" class="address" rows="2"></textarea>

            <label for="address">AMOUNT</label>
            <div class="amountbox">
                <input type="text" class="whole" novalue="0" maxlength="7" filter="0~9">
                <input type="text" class="divider" readonly value=".">
                <input type="text" class="decimal" novalue="000000" maxlength="6" filter="0~9">
            </div>

            <input class="continue" type="submit" value="SEND PAYMENT">

        </div>
    </section>

    <section page="receive">
        <div class="middle">

            <label for="address">YOUR ADDRESS</label>
            <textarea class="address" rows="2" spellcheck="false"></textarea>

        </div>
    </section>

    <section page="representative">
        <div class="middle">

            <label for="address">ADDRESS</label>
            <textarea spellcheck="false" class="address" rows="2"></textarea>

            <input class="continue" type="submit" value="CHANGE REPRESENTATIVE">

        </div>
    </section>

    <section page="history">
        <div class="middle">

            <label>BALANCE</label>
            <div class="box fullamoutbox">
                <div class="item fullamount"></div>
            </div>

            <label>TRANSACTIONS</label>
            <div class="box txbox">
                <div class="item">
                    No transactions were found. :(
                </div>
            </div>

        </div>
    </section>


</main>`)
