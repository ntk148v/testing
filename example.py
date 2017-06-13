#!/usr/bin/python

"""
    example.py

    Python Daemon example. Creates a simple web sever and runs as a daemon
    
    Usage:
        1. Run from a terminal / command line: "python example.py start".
        2. Notice that the terminal is available immediately after the app has been daemonized.
        3. From the terminal run "python example.py status". The PID of the daemonized process should be displayed.
        4. Open a browser (on the same computer if the IP is localhost) and go to http://127.0.0.1:8001 (or your defeinfe IP adddress and Port).
        5. Answer the important question post on the simple website.
        6. From the terminal run "python example.py stop". The daemon should now be stopped.
        7. From the terminal run "python example.py status". The staus should show that the daemon is not running.
    
    Requires:
        Python 2.x
        BaseHHTPServer, SimpleHTTPServer modules from the Python standard library

    v 1.0
    23 Jan 2013
"""

# Configuration
# Full path names are required because the daemon has no knowledge of the current directory/folder
SETTINGS_LINUX = {
    'APP' :     'mydaemon_example',
    'WWW' :     '/home/Daemon/example.html',
    'PIDFILE' : '/home/Daemon/mydaemon_example.pid',
    'LOG' :     '/home/Daemon/mydaemon_example.log'
}
SETTINGS_WINDOWS = {
    # Required on if runing under Windows (supports in the foreground only)
    'APP' :     'mydaemon_example',
    'WWW' :     'H:\\Daemon\\example.html',
    'PIDFILE' : 'H:\\Daemon\\mydaemon_example.pid',
    'LOG' :     'H:\\Daemon\\mydaemon_example.log'}
SETTINGS_OTHER = {
    # Other Operating Systems/Platforms
    'APP' :     'mydaemon_example',
    'WWW' :     'index.html',
    'PIDFILE' : 'mydaemon_example.pid',
    'LOG' :     'mydaemon_example.log'}

# Web server settings
IP = '127.0.0.1'    # localhost is 127.0.0.1 or enter the IP address of your computer
PORT = 8002

# Std Python Modules
import sys, types, os, platform
from BaseHTTPServer import HTTPServer
from SimpleHTTPServer import SimpleHTTPRequestHandler

# App modules
from daemon import DaemonClass

# Globals for this module
__version__ = '1.0'

class MyApp(DaemonClass):
    """
        Subclasses the Daemon Class from the daemon module.
        Overrides relevant methods.
        Configures a basic web server.
    """

    def __init__(self, settings):
        # Initialise the parent class
        DaemonClass.__init__(self, settings)
        self.html = ""

    def run(self):
        """
            Runs after the process has been demonised.
            Overrides the run method provided by the class (which does nothing much anyway).
            When this method exits, the application exits.
        """

        sys.stdout.write("(pid:%s) Running MainLoop...\n" % self.pid)
        self.signaltrapped = False

        sys.stdout.write('Starting Web server...\n')
        try:
            server_address = (IP, PORT)
            webserver = HTTPServer(server_address, myHTTPHandler)
            sock = webserver.socket.getsockname()
            sys.stdout.write("Serving HTTP on " + str(sock[0]) + ", port " + str(sock[1]) + " ...")
            # Blocks here, serving HTTP GET requests (via myHTTPHandler), until a signal is received
            webserver.serve_forever()
        except:
            sys.stdout.write("Error starting HTTP server")

        # Interrupt received so quit gracefully
        sys.stdout.write('...Shutting down the daemon.\n')

        # Final tidy up - delete PID file etc
        self.on_exit()
        self.cleanup()
        sys.stdout.write("All Done.")
        # The daemon will close down now

    def on_interrupt(self):
        """
            Quit gracefully
        """
        sys.stdout.write("Interrupt received, so quitting gracefully...\n")

class myHTTPHandler(SimpleHTTPRequestHandler):
    """
        Handler for the GET requests
        Subclasses BaseHTTPRequestHandler
    """
    def do_GET(self):
        """
            Respond to HTTP GET requests. Overrides the Std Library's do_GET method
        """
        self.html = ''
        self.send_response(200)
        self.send_header('Content-type','text/html')
        self.end_headers()
        # Send the html message
        try:
            f = open(SETTINGS_LINUX['WWW'], 'r')
            self.html = f.read()
            f.close()
        except:
            self.html = "Cannot open website content:" % self.settings['www']
            sys.stderr.write(self.html)
        
        self.wfile.write(self.html)
        return

# End Class


if __name__ == '__main__':
    global app

    if platform.system() == 'Linux':
        SETTINGS = SETTINGS_LINUX
    elif platform.system() == 'Windows':
        SETTINGS = SETTINGS_WINDOWS
    else:
        SETTINGS = SETTINGS_OTHER

    print "----- %s (Release: %s) -----" % (SETTINGS['APP'], __version__)

    if len(sys.argv) == 2:
        if 'start' == sys.argv[1]:
            sys.stdout.write("Starting the app...\n")
            app = MyApp(SETTINGS)
            sys.stdout.write("Starting daemon mode...")
            app.start()
        elif 'stop' == sys.argv[1]:
            app = MyApp(SETTINGS)
            sys.stdout.write("Stopping the daemon...\n")
            app.stop()
        elif 'restart' == sys.argv[1]:
            app = MyApp(SETTINGS)
            sys.stdout.write("Restarting the daemon...")
            app.restart()
        elif 'status' == sys.argv[1]:
            app = MyApp(SETTINGS)
            app.status()
        else:
            print "usage: %s start|stop|restart/status" % sys.argv[0]
            sys.exit(2)
    else:
            print "Invalid command: %r" % ' '.join(sys.argv)
            print "usage: %s start|stop|restart|status" % sys.argv[0]
            sys.exit(2)

    print "...Done"

# end of module
