INSERT INTO cameras (id, camera_slug, camera_model, camera_make, camera_type, camera_owner, camera_description, camera_notes, created_at, updated_at, deleted_at) VALUES (1, 'unknown', 'Unknown', '', '', '', '', '', '2020-01-06 02:06:29', '2020-01-06 02:07:26', null);
INSERT INTO cameras (id, camera_slug, camera_model, camera_make, camera_type, camera_owner, camera_description, camera_notes, created_at, updated_at, deleted_at) VALUES (2, 'apple-iphone-se', 'iPhone SE', 'Apple', '', '', '', '', '2020-01-06 02:06:30', '2020-01-06 02:07:28', null);
INSERT INTO cameras (id, camera_slug, camera_model, camera_make, camera_type, camera_owner, camera_description, camera_notes, created_at, updated_at, deleted_at) VALUES (3, 'canon-eos-5d', 'EOS 5D', 'Canon', '', '', '', '', '2020-01-06 02:06:32', '2020-01-06 02:06:32', null);
INSERT INTO cameras (id, camera_slug, camera_model, camera_make, camera_type, camera_owner, camera_description, camera_notes, created_at, updated_at, deleted_at) VALUES (4, 'canon-eos-7d', 'EOS 7D', 'Canon', '', '', '', '', '2020-01-06 02:06:33', '2020-01-06 02:06:33', null);
INSERT INTO cameras (id, camera_slug, camera_model, camera_make, camera_type, camera_owner, camera_description, camera_notes, created_at, updated_at, deleted_at) VALUES (5, 'canon-eos-6d', 'EOS 6D', 'Canon', '', '', '', '', '2020-01-06 02:06:35', '2020-01-06 02:06:54', null);
INSERT INTO cameras (id, camera_slug, camera_model, camera_make, camera_type, camera_owner, camera_description, camera_notes, created_at, updated_at, deleted_at) VALUES (6, 'apple-iphone-6', 'iPhone 6', 'Apple', '', '', '', '', '2020-01-06 02:06:42', '2020-01-06 02:06:42', null);
INSERT INTO cameras (id, camera_slug, camera_model, camera_make, camera_type, camera_owner, camera_description, camera_notes, created_at, updated_at, deleted_at) VALUES (7, 'apple-iphone-7', 'iPhone 7', 'Apple', '', '', '', '', '2020-01-06 02:06:51', '2020-01-06 02:06:51', null);
INSERT INTO countries (id, country_slug, country_name, country_description, country_notes, country_photo_id) VALUES ('de', 'germany', 'Germany', 'Country Description', 'Country Notes', 0);
INSERT INTO albums (id, album_uuid, album_name, album_slug, album_favorite) VALUES ('2', '3', 'Christmas2030', 'christmas2030', 0);
INSERT INTO albums (id, album_uuid, cover_uuid, album_name, album_slug, album_favorite) VALUES ('1', '4', '654', 'Holiday2030', 'holiday-2030', 1);
INSERT INTO photos_albums (album_uuid, photo_uuid) VALUES ('4', '654');
INSERT INTO files (id, photo_id, photo_uuid, file_name, file_primary) VALUES ('1', '1', '654', 'exampleFileName.jpg', 1);
INSERT INTO photos (id, photo_uuid) VALUES ('1', '654');
INSERT INTO categories (label_id, category_id) VALUES ('1', '1');
INSERT INTO labels (id, label_name) VALUES ('1', 'flower');



