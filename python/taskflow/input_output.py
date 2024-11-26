from taskflow import engines
from taskflow.patterns import linear_flow as lf
from taskflow.patterns import unordered_flow as uf
from taskflow import task


# class TaskA(task.Task):
#     default_provides = 'a'

#     def execute(self):
#         print("Executing '%s'" % (self.name))
#         return 'a'


# class TaskB(task.Task):
#     default_provides = 'b'

#     def execute(self, a):
#         print("Executing '%s'" % (self.name))
#         print("Got input '%s'" % (a))
#         a = 'b'
#         return a


# print("Constructing...")
# wf = lf.Flow("pass-from-to")
# wf.add(TaskA('a'), TaskB('b'))

# print("Loading...")
# e = engines.load(wf)

# print("Running...")
# print(e.run())

# print("Done...")


class CatTalk(task.Task):
    def execute(self, meow):
        print(meow)
        return "cat"


class DogTalk(task.Task):
    default_provides = ["dog"]

    def execute(self, woof):
        print(woof)
        return {"dog": "dog"}


flo = lf.Flow("cat-dog")
flo.add(CatTalk(), DogTalk())
e = engines.load(flo, store={'meow': 'meow', 'woof': 'woof'})
print(e.run())
