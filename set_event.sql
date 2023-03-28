CREATE OR REPLACE PROCEDURE sh_ilo.set_event(IN p_id bigint, IN p_device bigint, IN p_object integer, IN p_event integer, IN p_message jsonb, IN p_info jsonb)
 LANGUAGE plpgsql
 SECURITY DEFINER
AS $procedure$
    begin
        if p_event > 0 then
            INSERT INTO sh_ilo.events_messages (created_id, object_id, device_id, mes_id, mes_time, mes_code, mes_status, mes_data, event_value, event_data)
            SELECT convert_to((p_message ->> '_id'), 'utf8'),
                    p_object,
                    p_device,
                    (p_message -> 'mes_id') :: bigint,
                    (p_message ->> 'mes_time') :: timestamp,
                    c.code,
                    (p_message -> 'status_info'),
                    p_message,
                    c.const,
                    p_info
            FROM sh_data.dir_constants c where c.code = p_event and c.class in ('DBMSGTYPE', 'DBLOGICTYPE');
        end if;

        INSERT INTO sh_ilo.data_device_event (id, object_id, device_id, created_id, mes_id, mes_time, event, data, event_data)
        VALUES (p_id,
                p_object,
                p_device,
                convert_to((p_message ->> '_id'), 'utf8'),
                (p_message -> 'mes_id') :: bigint,
                (p_message ->> 'mes_time') :: timestamp,
                p_event,
                p_message,
                p_info)
        ON CONFLICT (object_id, event)
            DO UPDATE SET id          = excluded.id,
                          created_id=excluded.created_id,
                          device_id=excluded.device_id,
                          mes_id=excluded.mes_id,
                          mes_time=excluded.mes_time,
                          data=excluded.data,
                          event_data  = excluded.event_data,
                          update_time = CURRENT_TIMESTAMP;

    end;

$procedure$
;