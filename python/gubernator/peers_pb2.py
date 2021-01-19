# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: peers.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import gubernator_pb2 as gubernator__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='peers.proto',
  package='pb.gubernator',
  syntax='proto3',
  serialized_options=b'Z\014.;gubernator\200\001\001',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x0bpeers.proto\x12\rpb.gubernator\x1a\x10gubernator.proto\"E\n\x14GetPeerRateLimitsReq\x12-\n\x08requests\x18\x01 \x03(\x0b\x32\x1b.pb.gubernator.RateLimitReq\"J\n\x15GetPeerRateLimitsResp\x12\x31\n\x0brate_limits\x18\x01 \x03(\x0b\x32\x1c.pb.gubernator.RateLimitResp\"H\n\x14UpdatePeerGlobalsReq\x12\x30\n\x07globals\x18\x01 \x03(\x0b\x32\x1f.pb.gubernator.UpdatePeerGlobal\"z\n\x10UpdatePeerGlobal\x12\x0b\n\x03key\x18\x01 \x01(\t\x12,\n\x06status\x18\x02 \x01(\x0b\x32\x1c.pb.gubernator.RateLimitResp\x12+\n\talgorithm\x18\x03 \x01(\x0e\x32\x18.pb.gubernator.Algorithm\"\x17\n\x15UpdatePeerGlobalsResp2\xcd\x01\n\x07PeersV1\x12`\n\x11GetPeerRateLimits\x12#.pb.gubernator.GetPeerRateLimitsReq\x1a$.pb.gubernator.GetPeerRateLimitsResp\"\x00\x12`\n\x11UpdatePeerGlobals\x12#.pb.gubernator.UpdatePeerGlobalsReq\x1a$.pb.gubernator.UpdatePeerGlobalsResp\"\x00\x42\x11Z\x0c.;gubernator\x80\x01\x01\x62\x06proto3'
  ,
  dependencies=[gubernator__pb2.DESCRIPTOR,])




_GETPEERRATELIMITSREQ = _descriptor.Descriptor(
  name='GetPeerRateLimitsReq',
  full_name='pb.gubernator.GetPeerRateLimitsReq',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='requests', full_name='pb.gubernator.GetPeerRateLimitsReq.requests', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=48,
  serialized_end=117,
)


_GETPEERRATELIMITSRESP = _descriptor.Descriptor(
  name='GetPeerRateLimitsResp',
  full_name='pb.gubernator.GetPeerRateLimitsResp',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='rate_limits', full_name='pb.gubernator.GetPeerRateLimitsResp.rate_limits', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=119,
  serialized_end=193,
)


_UPDATEPEERGLOBALSREQ = _descriptor.Descriptor(
  name='UpdatePeerGlobalsReq',
  full_name='pb.gubernator.UpdatePeerGlobalsReq',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='globals', full_name='pb.gubernator.UpdatePeerGlobalsReq.globals', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=195,
  serialized_end=267,
)


_UPDATEPEERGLOBAL = _descriptor.Descriptor(
  name='UpdatePeerGlobal',
  full_name='pb.gubernator.UpdatePeerGlobal',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='pb.gubernator.UpdatePeerGlobal.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='status', full_name='pb.gubernator.UpdatePeerGlobal.status', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='algorithm', full_name='pb.gubernator.UpdatePeerGlobal.algorithm', index=2,
      number=3, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=269,
  serialized_end=391,
)


_UPDATEPEERGLOBALSRESP = _descriptor.Descriptor(
  name='UpdatePeerGlobalsResp',
  full_name='pb.gubernator.UpdatePeerGlobalsResp',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=393,
  serialized_end=416,
)

_GETPEERRATELIMITSREQ.fields_by_name['requests'].message_type = gubernator__pb2._RATELIMITREQ
_GETPEERRATELIMITSRESP.fields_by_name['rate_limits'].message_type = gubernator__pb2._RATELIMITRESP
_UPDATEPEERGLOBALSREQ.fields_by_name['globals'].message_type = _UPDATEPEERGLOBAL
_UPDATEPEERGLOBAL.fields_by_name['status'].message_type = gubernator__pb2._RATELIMITRESP
_UPDATEPEERGLOBAL.fields_by_name['algorithm'].enum_type = gubernator__pb2._ALGORITHM
DESCRIPTOR.message_types_by_name['GetPeerRateLimitsReq'] = _GETPEERRATELIMITSREQ
DESCRIPTOR.message_types_by_name['GetPeerRateLimitsResp'] = _GETPEERRATELIMITSRESP
DESCRIPTOR.message_types_by_name['UpdatePeerGlobalsReq'] = _UPDATEPEERGLOBALSREQ
DESCRIPTOR.message_types_by_name['UpdatePeerGlobal'] = _UPDATEPEERGLOBAL
DESCRIPTOR.message_types_by_name['UpdatePeerGlobalsResp'] = _UPDATEPEERGLOBALSRESP
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

GetPeerRateLimitsReq = _reflection.GeneratedProtocolMessageType('GetPeerRateLimitsReq', (_message.Message,), {
  'DESCRIPTOR' : _GETPEERRATELIMITSREQ,
  '__module__' : 'peers_pb2'
  # @@protoc_insertion_point(class_scope:pb.gubernator.GetPeerRateLimitsReq)
  })
_sym_db.RegisterMessage(GetPeerRateLimitsReq)

GetPeerRateLimitsResp = _reflection.GeneratedProtocolMessageType('GetPeerRateLimitsResp', (_message.Message,), {
  'DESCRIPTOR' : _GETPEERRATELIMITSRESP,
  '__module__' : 'peers_pb2'
  # @@protoc_insertion_point(class_scope:pb.gubernator.GetPeerRateLimitsResp)
  })
_sym_db.RegisterMessage(GetPeerRateLimitsResp)

UpdatePeerGlobalsReq = _reflection.GeneratedProtocolMessageType('UpdatePeerGlobalsReq', (_message.Message,), {
  'DESCRIPTOR' : _UPDATEPEERGLOBALSREQ,
  '__module__' : 'peers_pb2'
  # @@protoc_insertion_point(class_scope:pb.gubernator.UpdatePeerGlobalsReq)
  })
_sym_db.RegisterMessage(UpdatePeerGlobalsReq)

UpdatePeerGlobal = _reflection.GeneratedProtocolMessageType('UpdatePeerGlobal', (_message.Message,), {
  'DESCRIPTOR' : _UPDATEPEERGLOBAL,
  '__module__' : 'peers_pb2'
  # @@protoc_insertion_point(class_scope:pb.gubernator.UpdatePeerGlobal)
  })
_sym_db.RegisterMessage(UpdatePeerGlobal)

UpdatePeerGlobalsResp = _reflection.GeneratedProtocolMessageType('UpdatePeerGlobalsResp', (_message.Message,), {
  'DESCRIPTOR' : _UPDATEPEERGLOBALSRESP,
  '__module__' : 'peers_pb2'
  # @@protoc_insertion_point(class_scope:pb.gubernator.UpdatePeerGlobalsResp)
  })
_sym_db.RegisterMessage(UpdatePeerGlobalsResp)


DESCRIPTOR._options = None

_PEERSV1 = _descriptor.ServiceDescriptor(
  name='PeersV1',
  full_name='pb.gubernator.PeersV1',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=419,
  serialized_end=624,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetPeerRateLimits',
    full_name='pb.gubernator.PeersV1.GetPeerRateLimits',
    index=0,
    containing_service=None,
    input_type=_GETPEERRATELIMITSREQ,
    output_type=_GETPEERRATELIMITSRESP,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
  _descriptor.MethodDescriptor(
    name='UpdatePeerGlobals',
    full_name='pb.gubernator.PeersV1.UpdatePeerGlobals',
    index=1,
    containing_service=None,
    input_type=_UPDATEPEERGLOBALSREQ,
    output_type=_UPDATEPEERGLOBALSRESP,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_PEERSV1)

DESCRIPTOR.services_by_name['PeersV1'] = _PEERSV1

# @@protoc_insertion_point(module_scope)
