from st2common.runners.base_action import Action


class EchoAction(Action):
    def run(self, message):
        print(message)

        if message == "working":
            # You can use self.logger for logging
            self.logger.info("Action successfully completed")
            return (True, message)
        self.logger.error("Action failed")
        return (False, message)
