package testdata;

import org.junit.Test;

import java.io.ByteArrayOutputStream;
import java.io.ByteArrayInputStream;
import java.math.BigInteger;
import java.nio.ByteBuffer;
import java.time.Instant;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.Map;
import java.util.Map.Entry;

import static java.nio.charset.StandardCharsets.UTF_8;
import static org.junit.Assert.assertArrayEquals;
import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNull;
import static org.junit.Assert.fail;


public class test {

	static Map<String, O> newGoldenCases() {
		Map<String, O> goldenCases = new LinkedHashMap<>();
		newCase(goldenCases, "7f");
		newCase(goldenCases, "007f").b = true;
		newCase(goldenCases, "01017f").u32 = 1;
		newCase(goldenCases, "01ff017f").u32 = 255;
		newCase(goldenCases, "01ffff037f").u32 = 65535;
		newCase(goldenCases, "81ffffffff7f").u32 = -1;
		newCase(goldenCases, "02017f").u64 = 1L;
		newCase(goldenCases, "02ff017f").u64 = 255L;
		newCase(goldenCases, "02ffff037f").u64 = 65535L;
		newCase(goldenCases, "02ffffffff0f7f").u64 = 4294967295L;
		newCase(goldenCases, "82ffffffffffffffff7f").u64 = -1L;
		newCase(goldenCases, "03017f").i32 = 1;
		newCase(goldenCases, "83017f").i32 = -1;
		newCase(goldenCases, "037f7f").i32 = Byte.MAX_VALUE;
		newCase(goldenCases, "8380017f").i32 = Byte.MIN_VALUE;
		newCase(goldenCases, "03ffff017f").i32 = Short.MAX_VALUE;
		newCase(goldenCases, "838080027f").i32 = Short.MIN_VALUE;
		newCase(goldenCases, "03ffffffff077f").i32 = Integer.MAX_VALUE;
		newCase(goldenCases, "8380808080087f").i32 = Integer.MIN_VALUE;
		newCase(goldenCases, "04017f").i64 = 1;
		newCase(goldenCases, "84017f").i64 = -1;
		newCase(goldenCases, "047f7f").i64 = Byte.MAX_VALUE;
		newCase(goldenCases, "8480017f").i64 = Byte.MIN_VALUE;
		newCase(goldenCases, "04ffff017f").i64 = Short.MAX_VALUE;
		newCase(goldenCases, "848080027f").i64 = Short.MIN_VALUE;
		newCase(goldenCases, "04ffffffff077f").i64 = Integer.MAX_VALUE;
		newCase(goldenCases, "8480808080087f").i64 = Integer.MIN_VALUE;
		newCase(goldenCases, "04ffffffffffffffff7f7f").i64 = Long.MAX_VALUE;
		newCase(goldenCases, "848080808080808080807f").i64 = Long.MIN_VALUE;
		newCase(goldenCases, "05000000017f").f32 = Float.MIN_VALUE;
		newCase(goldenCases, "057f7fffff7f").f32 = Float.MAX_VALUE;
		newCase(goldenCases, "057fc000007f").f32 = Float.NaN;
		newCase(goldenCases, "0600000000000000017f").f64 = Double.MIN_VALUE;
		newCase(goldenCases, "067fefffffffffffff7f").f64 = Double.MAX_VALUE;
		newCase(goldenCases, "067ff80000000000007f").f64 = Double.NaN;
		newCase(goldenCases, "0755ef312a2e5da4e77f").t = Instant.ofEpochSecond(1441739050L, 777888999);
		newCase(goldenCases, "870000000100000000000000007f").t = Instant.ofEpochSecond(1L << 32, 0);
		newCase(goldenCases, "87ffffffffffffffff2e5da4e77f").t = Instant.ofEpochSecond(-1L, 777888999);
		newCase(goldenCases, "87fffffff14f443f00000000007f").t = Instant.ofEpochSecond(-63094636800L, 0);
		newCase(goldenCases, "0801417f").s = "A";
		newCase(goldenCases, "080261007f").s = "a\0";
		newCase(goldenCases, "0809c280e0a080f09080807f").s = "\u0080\u0800\ud800\udc00";
		newCase(goldenCases, "08800120202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020202020207f").s = "                                                                                                                                ";
		newCase(goldenCases, "0901ff7f").a = new byte[]{-1};
		newCase(goldenCases, "090202007f").a = new byte[]{2, 0};
		newCase(goldenCases, "09c0010909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909090909097f").a = "\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t\t".getBytes(UTF_8);
		newCase(goldenCases, "0a7f7f").o = new O();
		O inner = new O();
		inner.b = true;
		newCase(goldenCases, "0a007f7f").o = inner;
		O element = new O();
		element.b = true;
		newCase(goldenCases, "0b01007f7f").os = new O[] {element};
		newCase(goldenCases, "0b027f7f7f").os = new O[] {new O(), new O()};
		newCase(goldenCases, "0c0300016101627f").ss = new String[] {"", "a", "b"};
		newCase(goldenCases, "0d0201000201027f").as = new byte[][]{new byte[]{0}, new byte[]{1, 2}};
		return goldenCases;
	}

	private static O newCase(Map<String, O> cases, String hex) {
		O o = new O();
		cases.put(hex, o);
		return o;
	}

	@Test
	public void marshal() {
		for (Entry<String, O> e : newGoldenCases().entrySet()) {
			try {
				byte[] buf = new byte[e.getKey().length() / 2];
				int n = e.getValue().marshal(buf, 0);
				assertEquals("serial", e.getKey(), toHex(buf));
				assertEquals("write index", n, buf.length);
			} catch (Exception ex) {
				ex.printStackTrace();
				fail("exception for serial " + e.getKey());
			}
		}
	}

	@Test
	public void unmarshal() {
		for (Entry<String, O> e : newGoldenCases().entrySet()) {
			try {
				O o = new O();
				byte[] serial = parseHex(e.getKey());
				int n = o.unmarshal(serial, 0);
				assertEquals(e.getKey(), e.getValue(), o);
				assertEquals("read index", n, serial.length);
			} catch (Exception ex) {
				ex.printStackTrace();
				fail("exception for serial " + e.getKey());
			}
		}
	}

	@Test
	public void streaming() throws Exception {
		ByteArrayOutputStream out = new ByteArrayOutputStream();

		byte[] buf = new byte[1];
		for (O o : newGoldenCases().values()) {
			buf = o.marshal(out, buf);
		}

		O.Unmarshaller unmarshaller = new O.Unmarshaller(new ByteArrayInputStream(out.toByteArray()), new byte[1]);
		for (O o : newGoldenCases().values()) {
			assertEquals(unmarshaller.next(), o);
		}
		assertNull(unmarshaller.next());
	}

	@Test
	public void marshalMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 2;
		try {
			O o = new O();
			o.u64 = 1;
			o.marshal(new byte[O.colferSizeMax], 0);
			fail("no marshal exception");
		} catch (IllegalStateException e) {
			assertEquals("marshal error", "colfer: testdata.o exceeds 2 bytes", e.getMessage());
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	@Test
	public void marshalTextMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 2;
		try {
			O o = new O();
			o.s = "AAA";
			o.marshal(new byte[6], 0);
			fail("no marshal exception");
		} catch (IllegalStateException e) {
			// Field message only when buffer is big enough. Otherwise it's: "serial exceeds 2 bytes".
			assertEquals("marshal error", "colfer: testdata.o.s size 3 exceeds 2 UTF-8 bytes", e.getMessage());
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	@Test
	public void marshalBinaryMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 2;
		try {
			O o = new O();
			o.a = new byte[]{0, 1, 2};
			o.marshal(new byte[O.colferSizeMax], 0);
			fail("no marshal exception");
		} catch (IllegalStateException e) {
			assertEquals("marshal error", "colfer: testdata.o.a size 3 exceeds 2 bytes", e.getMessage());
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	@Test
	public void marshalListMax() {
		int origMax = O.colferListMax;
		O.colferListMax = 9;
		try {
			O o = new O();
			o.os = new O[10];
			o.marshal(new byte[O.colferSizeMax], 0);
			fail("no marshal exception");
		} catch (IllegalStateException e) {
			assertEquals("marshal error", "colfer: testdata.o.os length 10 exceeds 9 elements", e.getMessage());
		} finally {
			O.colferListMax = origMax;
		}
	}

	@Test
	public void unmarshalMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 2;
		try {
			byte[] serial = parseHex("02017f");
			new O().unmarshal(serial, 0);
			fail("no unmarshal exception");
		} catch (SecurityException e) {
			assertEquals("unmarshal error", "colfer: testdata.o exceeds 2 bytes", e.getMessage());
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	@Test
	public void unmarshalTextMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 9;
		try {
			byte[] serial = parseHex("080a414141");
			new O().unmarshal(serial, 0);
			fail("no unmarshal exception");
		} catch (SecurityException e) {
			assertEquals("unmarshal error", "colfer: testdata.o.s size 10 exceeds 9 UTF-8 bytes", e.getMessage());
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	@Test
	public void unmarshalBinaryMax() {
		int origMax = O.colferSizeMax;
		O.colferSizeMax = 9;
		try {
			byte[] serial = parseHex("090a414141");
			new O().unmarshal(serial, 0);
			fail("no unmarshal exception");
		} catch (SecurityException e) {
			assertEquals("unmarshal error", "colfer: testdata.o.a size 10 exceeds 9 bytes", e.getMessage());
		} finally {
			O.colferSizeMax = origMax;
		}
	}

	@Test
	public void unmarshalListMax() {
		int origMax = O.colferListMax;
		O.colferListMax = 9;
		try {
			byte[] serial = parseHex("0b0a7f7f7f");
			new O().unmarshal(serial, 0);
			fail("no unmarshal exception");
		} catch (SecurityException e) {
			assertEquals("unmarshal error", "colfer: testdata.o.os length 10 exceeds 9 elements", e.getMessage());
		} finally {
			O.colferListMax = origMax;
		}
	}

	static String toHex(byte[] bytes) {
		String hex = new BigInteger(1, bytes).toString(16);
		while (bytes.length * 2 > hex.length())
			hex = "0" + hex;
		return hex;
	}

	static byte[] parseHex(String s) {
		int len = s.length();
		byte[] data = new byte[len / 2];
		for (int i = 0; i < len; i += 2) {
			int nibble0 = Character.digit(s.charAt(i), 16);
			int nibble1 = Character.digit(s.charAt(i + 1), 16);
			data[i / 2] = (byte) ((nibble0 << 4) + nibble1);
		}
		return data;
	}

	@Test
	public void identity() {
		Object[] a = newGoldenCases().values().toArray();
		Object[] b = newGoldenCases().values().toArray();
		assertArrayEquals("golden cases", a, b);
		assertEquals("golden cases hash", Arrays.hashCode(a), Arrays.hashCode(b));
	}

}
