TESLAMATE
===================================================================

P1
-------------------------------------------------------------------

SELECT sum(charge_energy_used) FROM charging_processes
WHERE start_date >= '2020-09-01'
AND start_date < '2020-10-01'
AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '13:00' and '23:00'
AND CAST(end_date at time zone 'gmt+2' AS time) BETWEEN '13:00' and '23:00'

P2
-------------------------------------------------------------------

SELECT sum(charge_energy_used) charge_energy_used FROM
(
    SELECT sum(charge_energy_used) charge_energy_used FROM charging_processes
    WHERE start_date >= '2020-09-01'
    AND start_date < '2020-10-01'
    AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '07:00' and '13:00'
    AND CAST(end_date at time zone 'gmt+2' AS time) BETWEEN '07:00' and '13:00'
    UNION
    SELECT sum(charge_energy_used) charge_energy_used FROM charging_processes
    WHERE start_date >= '2020-09-01'
    AND start_date < '2020-10-01'
    AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '23:00' and '23:59'
    AND CAST(end_date at time zone 'gmt+2' AS time) BETWEEN '23:00' and '23:59'
    UNION
    SELECT sum(charge_energy_used) charge_energy_used FROM charging_processes
    WHERE start_date >= '2020-09-01'
    AND start_date < '2020-10-01'
    AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '00:00' and '01:00'
    AND CAST(end_date at time zone 'gmt+2' AS time) BETWEEN '00:00' and '01:00'
)s

P3
-------------------------------------------------------------------

SELECT sum(charge_energy_used) FROM charging_processes
WHERE start_date >= '2020-09-01'
AND start_date < '2020-10-01'
AND CAST(start_date at time zone 'gmt+2' AS time) BETWEEN '01:00' and '07:00'
AND CAST(end_date at time zone 'gmt+2' AS time) BETWEEN '01:00' and '07:00'