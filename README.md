# Magasin Advicer
Created by Jedikrigeren for simply-chocolate

Purpose of this script is to send an "advice" to magasin when we are sending them goods. The Advice is a CSV file that is uploadet to their FTP server.

Further development:
  - Create a SAP field on Business Partner level where you choose which magasin house code is corresponds to the Business Partner (10, 15, 20, 25, 30, 40, 50, 60). This way Magasin can add new shops without the need to alter the script itself.
  - Create a SAP field on Document level that can be used to mark a StockTransfer as "Adviced". `U_CCF_Maga_Adviseret: Y ? N`

There are currently no known errors.

The script send a message to teams every time it has finished, so if you do not recieve a daily message on the teams channel, it means the script is not running.
