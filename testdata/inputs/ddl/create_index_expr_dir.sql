create index VenuesByCity on Venues (
  (json_value(VenueData.address.city)) desc
)
