package txt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDateTime(t *testing.T) {
	t.Run("2016: :     :  :  ", func(t *testing.T) {
		result := DateTime("2016: :     :  :  ", "")
		assert.Equal(t, "2016-01-01 00:00:00 +0000 UTC", result.String())
	})
	t.Run("2016:  :__    :  :  ", func(t *testing.T) {
		result := DateTime("2016:  :__   :  :  ", "")
		assert.Equal(t, "2016-01-01 00:00:00 +0000 UTC", result.String())
	})
	t.Run("2016:06:28   :  :??", func(t *testing.T) {
		result := DateTime("2016:06:28   :  :??", "")
		assert.Equal(t, "2016-06-28 00:00:00 +0000 UTC", result.String())
	})
	t.Run("2016:06:28 09:45:49", func(t *testing.T) {
		result := DateTime("2016:06:28 09:45:49", "")
		assert.Equal(t, "2016-06-28 09:45:49 +0000 UTC", result.String())
	})
	t.Run("2016:06:28 09:45:49+10:00", func(t *testing.T) {
		result := DateTime("2016:06:28 09:45:49+10:00", "")
		assert.Equal(t, "2016-06-28 09:45:49 +1000 UTC+10:00", result.String())
	})
	t.Run("2016:06:28   :  :", func(t *testing.T) {
		result := DateTime("2016:06:28   :  :", "")
		assert.Equal(t, "2016-06-28 00:00:00 +0000 UTC", result.String())
	})
	t.Run("2016/06/28T09-45:49", func(t *testing.T) {
		result := DateTime("2016/06/28T09-45:49", "")
		assert.Equal(t, "2016-06-28 09:45:49 +0000 UTC", result.String())
	})
	t.Run("2016:06:28T09:45:49Z", func(t *testing.T) {
		result := DateTime("2016:06:28T09:45:49Z", "")
		assert.Equal(t, "2016-06-28 09:45:49 +0000 UTC", result.String())
	})
	t.Run("2016:06:28T09:45:  Z", func(t *testing.T) {
		result := DateTime("2016:06:28T09:45:  Z", "")
		assert.Equal(t, "2016-06-28 09:45:00 +0000 UTC", result.String())
	})
	t.Run("2016:06:28T09:45:  ", func(t *testing.T) {
		result := DateTime("2016:06:28T09:45:  ", "")
		assert.Equal(t, "2016-06-28 09:45:00 +0000 UTC", result.String())
	})
	t.Run("2016:06:28T09:45:  ZABC", func(t *testing.T) {
		result := DateTime("2016:06:28T09:45:  ZABC", "")
		assert.Equal(t, "2016-06-28 09:45:00 +0000 UTC", result.String())
	})
	t.Run("2016:06:28T09:45:  ABC", func(t *testing.T) {
		result := DateTime("2016:06:28T09:45:  ABC", "")
		assert.Equal(t, "2016-06-28 09:45:00 +0000 UTC", result.String())
	})
	t.Run("2016:06:28 09:45:49+10:00ABC", func(t *testing.T) {
		result := DateTime("2016:06:28 09:45:49+10:00ABC", "")
		assert.Equal(t, "2016-06-28 09:45:49 +1000 UTC+10:00", result.String())
	})
	t.Run("  2016:06:28 09:45:49-01:30ABC", func(t *testing.T) {
		result := DateTime("  2016:06:28 09:45:49-01:30ABC", "")
		assert.Equal(t, "2016-06-28 09:45:49 -0130 UTC-01:30", result.String())
	})
	t.Run("2016:06:28 09:45:49-0130", func(t *testing.T) {
		result := DateTime("2016:06:28 09:45:49-0130", "")
		assert.Equal(t, "2016-06-28 09:45:49 -0130 UTC-01:30", result.String())
	})
	t.Run("UTC/016:06:28 09:45:49-0130", func(t *testing.T) {
		result := DateTime("2016:06:28 09:45:49-0130", "UTC")
		assert.Equal(t, "2016-06-28 11:15:49 +0000 UTC", result.String())
	})
	t.Run("UTC/016:06:28 09:45:49-0130", func(t *testing.T) {
		result := DateTime("2016:06:28 09:45:49.0130", "UTC")
		assert.Equal(t, "2016-06-28 09:45:49.013 +0000 UTC", result.String())
	})
	t.Run("2012:08:08 22:07:18", func(t *testing.T) {
		result := DateTime("2012:08:08 22:07:18", "")
		assert.Equal(t, "2012-08-08 22:07:18 +0000 UTC", result.String())
	})
	t.Run("2020-01-30_09-57-18", func(t *testing.T) {
		result := DateTime("2020-01-30_09-57-18", "")
		assert.Equal(t, "2020-01-30 09:57:18 +0000 UTC", result.String())
	})
	t.Run("EuropeBerlin/2016:06:28 09:45:49+10:00ABC", func(t *testing.T) {
		result := DateTime("2016:06:28 09:45:49+10:00ABC", "Europe/Berlin")
		assert.Equal(t, "2016-06-28 01:45:49 +0200 CEST", result.String())
	})
	t.Run("EuropeBerlin/  2016:06:28 09:45:49-01:30ABC", func(t *testing.T) {
		result := DateTime("  2016:06:28 09:45:49-01:30ABC", "Europe/Berlin")
		assert.Equal(t, "2016-06-28 13:15:49 +0200 CEST", result.String())
	})
	t.Run("EuropeBerlin/2012:08:08 22:07:18", func(t *testing.T) {
		result := DateTime("2012:08:08 22:07:18", "Europe/Berlin")
		assert.Equal(t, "2012-08-08 22:07:18 +0200 CEST", result.String())
	})
	t.Run("EuropeBerlin/2020-01-30_09-57-18", func(t *testing.T) {
		result := DateTime("2020-01-30_09-57-18", "Europe/Berlin")
		assert.Equal(t, "2020-01-30 09:57:18 +0100 CET", result.String())
	})
	t.Run("EuropeBerlin/2020-10-17-48-24.950488", func(t *testing.T) {
		result := DateTime("2020:10:17 17:48:24.9508123", "UTC")
		assert.Equal(t, "2020-10-17 17:48:24.9508123 +0000 UTC", result.UTC().String())
		assert.Equal(t, "2020-10-17 17:48:24.9508123", result.Format("2006-01-02 15:04:05.999999999"))
	})
}

func TestDateFromFilePath(t *testing.T) {
	t.Run("2016/08/18 iPhone/WRNI2074.jpg", func(t *testing.T) {
		result := DateFromFilePath("2016/08/18 iPhone/WRNI2074.jpg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2016-08-18 00:00:00 +0000 UTC", result.String())
	})

	t.Run("2016/08/18 iPhone/OZBJ8443.jpg", func(t *testing.T) {
		result := DateFromFilePath("2016/08/18 iPhone/OZBJ8443.jpg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2016-08-18 00:00:00 +0000 UTC", result.String())
	})

	t.Run("2018/04 - April/2018-04-12 19:24:49.gif", func(t *testing.T) {
		result := DateFromFilePath("2018/04 - April/2018-04-12 19:24:49.gif")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2018-04-12 19:24:49 +0000 UTC", result.String())
	})

	t.Run("2018", func(t *testing.T) {
		result := DateFromFilePath("2018")
		assert.True(t, result.IsZero())
	})

	t.Run("2018-04-12 19/24/49.gif", func(t *testing.T) {
		result := DateFromFilePath("2018-04-12 19/24/49.gif")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2018-04-12 19:24:49 +0000 UTC", result.String())
	})

	t.Run("/2020/1212/20130518_142022_3D657EBD.jpg", func(t *testing.T) {
		result := DateFromFilePath("/2020/1212/20130518_142022_3D657EBD.jpg")
		//assert.False(t, result.IsZero())
		assert.True(t, result.IsZero())
	})

	t.Run("20130518_142022_3D657EBD.jpg", func(t *testing.T) {
		result := DateFromFilePath("20130518_142022_3D657EBD.jpg")
		//assert.False(t, result.IsZero())
		assert.True(t, result.IsZero())
	})

	t.Run("telegram_2020_01_30_09_57_18.jpg", func(t *testing.T) {
		result := DateFromFilePath("telegram_2020_01_30_09_57_18.jpg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2020-01-30 09:57:18 +0000 UTC", result.String())
	})

	t.Run("Screenshot 2019_05_21 at 10.45.52.png", func(t *testing.T) {
		result := DateFromFilePath("Screenshot 2019_05_21 at 10.45.52.png")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-05-21 10:45:52 +0000 UTC", result.String())
	})

	t.Run("telegram_2020-01-30_09-57-18.jpg", func(t *testing.T) {
		result := DateFromFilePath("telegram_2020-01-30_09-57-18.jpg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2020-01-30 09:57:18 +0000 UTC", result.String())
	})

	t.Run("Screenshot 2019-05-21 at 10.45.52.png", func(t *testing.T) {
		result := DateFromFilePath("Screenshot 2019-05-21 at 10.45.52.png")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-05-21 10:45:52 +0000 UTC", result.String())
	})

	t.Run("telegram_2020-01-30_09-18.jpg", func(t *testing.T) {
		result := DateFromFilePath("telegram_2020-01-30_09-18.jpg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2020-01-30 00:00:00 +0000 UTC", result.String())
	})

	t.Run("Screenshot 2019-05-21 at 10545.52.png", func(t *testing.T) {
		result := DateFromFilePath("Screenshot 2019-05-21 at 10545.52.png")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-05-21 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/2019-05-21/file2314.JPG", func(t *testing.T) {
		result := DateFromFilePath("/2019-05-21/file2314.JPG")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-05-21 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/2019.05.21", func(t *testing.T) {
		result := DateFromFilePath("/2019.05.21")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/05.21.2019", func(t *testing.T) {
		result := DateFromFilePath("/05.21.2019")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/21.05.2019", func(t *testing.T) {
		result := DateFromFilePath("/21.05.2019")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("05/21/2019", func(t *testing.T) {
		result := DateFromFilePath("05/21/2019")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("2019-07-23", func(t *testing.T) {
		result := DateFromFilePath("2019-07-23")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-07-23 00:00:00 +0000 UTC", result.String())
	})

	t.Run("Photos/2015-01-14", func(t *testing.T) {
		result := DateFromFilePath("Photos/2015-01-14")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2015-01-14 00:00:00 +0000 UTC", result.String())
	})

	t.Run("21/05/2019", func(t *testing.T) {
		result := DateFromFilePath("21/05/2019")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("2019/05/21", func(t *testing.T) {
		result := DateFromFilePath("2019/05/21")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-05-21 00:00:00 +0000 UTC", result.String())
	})

	t.Run("2019/05/2145", func(t *testing.T) {
		result := DateFromFilePath("2019/05/2145")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/05/21/2019", func(t *testing.T) {
		result := DateFromFilePath("/05/21/2019")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/21/05/2019", func(t *testing.T) {
		result := DateFromFilePath("/21/05/2019")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/2019/05/21.jpeg", func(t *testing.T) {
		result := DateFromFilePath("/2019/05/21.jpeg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-05-21 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/2019/05/21/foo.txt", func(t *testing.T) {
		result := DateFromFilePath("/2019/05/21/foo.txt")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-05-21 00:00:00 +0000 UTC", result.String())
	})

	t.Run("2019/21/05", func(t *testing.T) {
		result := DateFromFilePath("2019/21/05")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/2019/05/21/foo.jpg", func(t *testing.T) {
		result := DateFromFilePath("/2019/05/21/foo.jpg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-05-21 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/2019/21/05/foo.jpg", func(t *testing.T) {
		result := DateFromFilePath("/2019/21/05/foo.jpg")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/2019/5/foo.jpg", func(t *testing.T) {
		result := DateFromFilePath("/2019/5/foo.jpg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-05-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/2019/1/3/foo.jpg", func(t *testing.T) {
		result := DateFromFilePath("/2019/1/3/foo.jpg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-01-03 00:00:00 +0000 UTC", result.String())
	})

	t.Run("/1989/1/3/foo.jpg", func(t *testing.T) {
		result := DateFromFilePath("/1989/1/3/foo.jpg")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("545452019/1/3/foo.jpg", func(t *testing.T) {
		result := DateFromFilePath("/2019/1/3/foo.jpg")
		assert.False(t, result.IsZero())
		assert.Equal(t, "2019-01-03 00:00:00 +0000 UTC", result.String())
	})

	t.Run("fo.jpg", func(t *testing.T) {
		result := DateFromFilePath("fo.jpg")
		assert.True(t, result.IsZero())
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("n >6", func(t *testing.T) {
		result := DateFromFilePath("2020-01-30_09-87-18-23.jpg")
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})

	t.Run("year < yearmin", func(t *testing.T) {
		result := DateFromFilePath("1020-01-30_09-57-18.jpg")
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})
	t.Run("hour > hourmax", func(t *testing.T) {
		result := DateFromFilePath("2020-01-30_25-57-18.jpg")
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})
	t.Run("invalid days", func(t *testing.T) {
		result := DateFromFilePath("2020-01-00.jpg")
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})
	t.Run("IMG-20191120-WA0001.jpg", func(t *testing.T) {
		result := DateFromFilePath("IMG-20191120-WA0001.jpg")
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})
	t.Run("VID-20191120-WA0001.jpg", func(t *testing.T) {
		result := DateFromFilePath("VID-20191120-WA0001.jpg")
		assert.Equal(t, "0001-01-01 00:00:00 +0000 UTC", result.String())
	})
}

func TestIsTime(t *testing.T) {
	t.Run("/2020/1212/20130518_142022_3D657EBD.jpg", func(t *testing.T) {
		assert.False(t, IsTime("/2020/1212/20130518_142022_3D657EBD.jpg"))
	})

	t.Run("telegram_2020_01_30_09_57_18.jpg", func(t *testing.T) {
		assert.False(t, IsTime("telegram_2020_01_30_09_57_18.jpg"))
	})

	t.Run("", func(t *testing.T) {
		assert.False(t, IsTime(""))
	})

	t.Run("Screenshot 2019_05_21 at 10.45.52.png", func(t *testing.T) {
		assert.False(t, IsTime("Screenshot 2019_05_21 at 10.45.52.png"))
	})

	t.Run("telegram_2020-01-30_09-57-18.jpg", func(t *testing.T) {
		assert.False(t, IsTime("telegram_2020-01-30_09-57-18.jpg"))
	})

	t.Run("2013-05-18", func(t *testing.T) {
		assert.True(t, IsTime("2013-05-18"))
	})

	t.Run("2013-05-18 12:01:01", func(t *testing.T) {
		assert.True(t, IsTime("2013-05-18 12:01:01"))
	})

	t.Run("20130518_142022", func(t *testing.T) {
		assert.True(t, IsTime("20130518_142022"))
	})

	t.Run("2020_01_30_09_57_18", func(t *testing.T) {
		assert.True(t, IsTime("2020_01_30_09_57_18"))
	})

	t.Run("2019_05_21 at 10.45.52", func(t *testing.T) {
		assert.True(t, IsTime("2019_05_21 at 10.45.52"))
	})

	t.Run("2020-01-30_09-57-18", func(t *testing.T) {
		assert.True(t, IsTime("2020-01-30_09-57-18"))
	})
}

func TestYear(t *testing.T) {
	t.Run("London 2002", func(t *testing.T) {
		result := Year("/2002/London 81/")
		assert.Equal(t, 2002, result)
	})

	t.Run("San Francisco 2019", func(t *testing.T) {
		result := Year("San Francisco 2019")
		assert.Equal(t, 2019, result)
	})

	t.Run("string with no number", func(t *testing.T) {
		result := Year("Born in the U.S.A. is a song written and performed by Bruce Springsteen...")
		assert.Equal(t, 0, result)
	})

	t.Run("file name", func(t *testing.T) {
		result := Year("/share/photos/243546/2003/01/myfile.jpg")
		assert.Equal(t, 2003, result)
	})

	t.Run("path", func(t *testing.T) {
		result := Year("/root/1981/London 2005")
		assert.Equal(t, 2005, result)
	})

	t.Run("empty string", func(t *testing.T) {
		result := Year("")
		assert.Equal(t, 0, result)
	})
}
