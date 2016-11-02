Sensor for [wind vane PEET BROS](http://www.peetbros.com/shop/item.aspx?itemid=137) for a raspberry pi (zero). The sensor is easy to repair.

Connect the two pairs to the GPIO 0 and GPIO 1 through pull ups and pull downs,
and watch the result (console and http)

GPIO 1 to lower reed switch (wind vane) and GPIO 0 to upper reed switch (wind
direction sensor)

schematic (for one pair, reed switch inside the wind vane):

       ,----- . . . --------------- 3V3
       |
    reed sw
       |
       '----- . . . --+--[R=1K0]--- GPIO
                      |
                      '--[R=10K]--- GND
