from oslo_versionedobjects import base
from oslo_versionedobjects import fields as obj_fields


# Ensure that we always register our object with an object registry,
# so that it can be deserialized from its primitive form.
@base.VersionedObjectRegistry.register
class TestObject(base.VersionedObject):
    """Simple test class with some data about it"""

    VERSION = '1.0'

    OBJ_PROJECT_NAMESPACE = 'test'

    fields = {
        'name': obj_fields.StringField(),
    }

    def __setattr__(self, key, value):
        if not (key[0:5] == '_obj_'
                or key[0:7] == '_change'
                or key == '_context'
                or key in list(self.fields)
                or key == 'FIELDS'
                or key == 'VERSION'
                or key == 'fields'):
            raise AttributeError(
                "Designate object '%(type)s' has no"
                "attribute '%(name)s'" % {
                    'type': self.obj_name(),
                    'name': key,
                })
        super(TestObject, self).__setattr__(key, value)


test_ins = TestObject(name='kien')
print('- The __str__() output of this new object: %s' % test_ins)
print('- The name field of the objects: %s' % test_ins.name)
test_ins_prim = test_ins.obj_to_primitive()
# Convert object to primitive
print('- Primitive representation of this object: %s' % test_ins_prim)

# Now convert the primitive back to an object
test_ins = TestObject.obj_from_primitive(test_ins_prim)

test_ins.obj_reset_changes()
print('- The __str__() output of this new (reconstructed) object: %s' % test_ins)

# Mutating a field and showing what changed
test_ins.name = 'another_name'
print('- After name change, the set of fields that have been mutated is: %s' % test_ins.obj_what_changed())
print('- Return object\'s name: %s' % test_ins.obj_name())

print('Dang Van Dai')
