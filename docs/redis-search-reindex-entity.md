# Reindex one ORM entity

The Redis Search Indexes dev-panel page provides **Reindex ID** for ORM-managed
indexes. Select an index, enter its entity's positive integer ID, and submit.
The action uses the existing dev-panel authentication.

```
POST /dev/redis-search/reindex-entity/:pool/:index/
Content-Type: application/json

{"ID":"21019"}
```

IDs are decimal strings so the full uint64 range survives JavaScript and JSON.
Pools must be configured in `App.RedisPools.Search`. The requested index must
already exist and correspond to a registered ORM entity in that pool. The list
endpoint exposes `Pool` and `CanReindexEntity` for each index.

Hitrix calls the ORM's registered database indexer starting at `ID - 1`, using a
restricted pusher that queues and flushes writes only for the requested key.
The existing ORM owns field mapping, namespaces and soft-delete filtering.
The indexer reads one standard page; no full-index rebuild is started and no
other document is updated. MySQL is not modified. An ID absent from that page,
including a row excluded by the ORM's soft-delete rules, returns HTTP 404 without
writing a document; it never silently reindexes the next ID.

After success, `ClearCacheByIDs` clears the selected entity's Redis cache and
local cache in the handling engine. It does not invalidate cached query results
or local caches in other processes. Reindexing does not change the source row:
repair incorrect database values separately before using this action.

Custom indexes are excluded because their document keys and ID-to-source lookup
have no generic ORM contract. Existing full-index reindex actions are unchanged.
No ORM source or dependency version change is required.
