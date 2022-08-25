from rocketry import Rocketry
from rocketry.conds import daily

app = Rocketry()

@app.task(daily)
def do_daily():
    print("Hello!")

#  @app.task(every("10 seconds"))
#  def do_continuously():
    #  ...

#  @app.task(daily.after("07:00"))
#  def do_daily_after_seven():
    #  ...

#  @app.task(hourly & time_of_day.between("22:00", "06:00"))
#  def do_hourly_at_night():
    #  ...

#  @app.task((weekly.on("Monday") | weekly.on("Saturday")) & time_of_day.after("10:00"))
#  def do_twice_a_week_after_ten():
    #  ...

if __name__ == '__main__':
    app.run()
